package repository

import (
	"go_project2/internal/app/ds"
	"time"
)

// GetOrCreateDraft возвращает черновик для пользователя (creatorID), если нет – создаёт
func (r *Repository) GetOrCreateDraft(creatorID uint) (*ds.Calculation, error) {
	var calc ds.Calculation
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&calc).Error
	if err == nil {
		return &calc, nil
	}
	// не найдено – создаём новый черновик
	calc = ds.Calculation{
		Status:    "черновик",
		CreatedAt: time.Now(),
		CreatorID: creatorID,
	}
	err = r.db.Create(&calc).Error
	if err != nil {
		return nil, err
	}
	return &calc, nil
}

// AddElectrolyteToCalculation добавляет раствор в заявку (черновик)
// Возвращает ошибку, если такой раствор уже есть в заявке (unique constraint)
func (r *Repository) AddElectrolyteToCalculation(calcID, electrolyteID uint, volume int, comment string) error {
	item := ds.CalculationItem{
		CalculationID: calcID,
		ElectrolyteID: electrolyteID,
		Volume:        volume,
		Comment:       comment,
	}
	return r.db.Create(&item).Error
}

// GetCalculationByID загружает расчёт со всеми элементами и данными растворов
func (r *Repository) GetCalculationByID(id int) (*ds.Calculation, []ds.CalculationItem, error) {
	var calc ds.Calculation
	err := r.db.First(&calc, id).Error
	if err != nil {
		return nil, nil, err
	}
	var items []ds.CalculationItem
	err = r.db.Preload("Electrolyte").Where("calculation_id = ?", id).Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	return &calc, items, nil
}

// GetCartCount возвращает количество позиций в черновике пользователя
func (r *Repository) GetCartCount(creatorID uint) int64 {
	var calc ds.Calculation
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&calc).Error
	if err != nil {
		return 0
	}
	var count int64
	r.db.Model(&ds.CalculationItem{}).Where("calculation_id = ?", calc.ID).Count(&count)
	return count
}

// DeleteCalculation логически удаляет заявку (меняет статус на "удалён") через raw SQL
func (r *Repository) DeleteCalculation(calcID uint) error {
	// Используем SQL UPDATE без ORM
	return r.db.Exec("UPDATE calculations SET status = 'удалён' WHERE id = ?", calcID).Error
}