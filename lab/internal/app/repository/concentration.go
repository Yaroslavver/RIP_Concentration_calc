package repository

import (
	"fmt"
	"lab/internal/app/ds"
	"time"
)

func (r *Repository) GetOrCreateDraft(creatorID uint) (*ds.Concentration, error) {
	var conc ds.Concentration
	
	// Ищем существующий черновик
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&conc).Error
	if err == nil {
		return &conc, nil
	}
	
	// Если черновик не найден, создаём новый
	// ВАЖНО: не указываем ID, пусть БД сама его сгенерирует
	conc = ds.Concentration{
		Status:      "черновик",
		CreatedAt:   time.Now(),
		CreatorID:   creatorID,
		Description: "",
		Result:      "",
	}
	
	// Создаём запись без указания ID
	err = r.db.Omit("ID").Create(&conc).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create draft: %w", err)
	}
	
	// Загружаем созданную запись, чтобы получить сгенерированный ID
	err = r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&conc).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load created draft: %w", err)
	}
	
	return &conc, nil
}

func (r *Repository) GetCartInfo(creatorID uint) (draftID uint, count int64, err error) {
	var conc ds.Concentration
	err = r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&conc).Error
	if err != nil {
		return 0, 0, nil
	}
	var cnt int64
	r.db.Model(&ds.ConcentrationItem{}).Where("concentration_id = ?", conc.ID).Count(&cnt)
	return conc.ID, cnt, nil
}

func (r *Repository) AddItemToDraft(creatorID, electrolyteID uint, volume int, comment string) error {
	draft, err := r.GetOrCreateDraft(creatorID)
	if err != nil {
		return err
	}
	var existing ds.ConcentrationItem
	err = r.db.Where("concentration_id = ? AND electrolyte_id = ?", draft.ID, electrolyteID).First(&existing).Error
	if err == nil {
		return fmt.Errorf("electrolyte already in draft")
	}
	item := ds.ConcentrationItem{
		ConcentrationID: draft.ID,
		ElectrolyteID:   electrolyteID,
		Volume:          volume,
		Comment:         comment,
	}
	return r.db.Create(&item).Error
}

func (r *Repository) UpdateItem(creatorID, itemID uint, volume int, comment string) error {
	var item ds.ConcentrationItem
	if err := r.db.First(&item, itemID).Error; err != nil {
		return err
	}
	var conc ds.Concentration
	if err := r.db.First(&conc, item.ConcentrationID).Error; err != nil {
		return err
	}
	if conc.CreatorID != creatorID {
		return fmt.Errorf("not your item")
	}
	return r.db.Model(&item).Updates(map[string]interface{}{
		"volume":  volume,
		"comment": comment,
	}).Error
}

func (r *Repository) DeleteItem(creatorID, itemID uint) error {
	var item ds.ConcentrationItem
	if err := r.db.First(&item, itemID).Error; err != nil {
		return err
	}
	var conc ds.Concentration
	if err := r.db.First(&conc, item.ConcentrationID).Error; err != nil {
		return err
	}
	if conc.CreatorID != creatorID {
		return fmt.Errorf("not your item")
	}
	return r.db.Delete(&item).Error
}

func (r *Repository) GetUserConcentrations(creatorID uint, status string, fromDate, toDate *time.Time) ([]ds.Concentration, error) {
	query := r.db.Where("creator_id = ? AND status NOT IN (?, ?)", creatorID, "черновик", "удалён")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != nil {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", toDate)
	}
	var concs []ds.Concentration
	err := query.Find(&concs).Error
	return concs, err
}

func (r *Repository) GetAllConcentrations(status string, fromDate, toDate *time.Time) ([]ds.Concentration, error) {
	query := r.db.Where("status NOT IN (?, ?)", "черновик", "удалён")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != nil {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", toDate)
	}
	var concs []ds.Concentration
	err := query.Preload("Creator").Find(&concs).Error
	return concs, err
}

func (r *Repository) GetConcentrationByID(id uint) (*ds.Concentration, []ds.ConcentrationItem, error) {
	var conc ds.Concentration
	err := r.db.First(&conc, id).Error
	if err != nil {
		return nil, nil, err
	}
	var items []ds.ConcentrationItem
	err = r.db.Preload("Electrolyte").Where("concentration_id = ?", id).Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	return &conc, items, nil
}

func (r *Repository) UpdateConcentration(conc *ds.Concentration) error {
	return r.db.Model(conc).Updates(map[string]interface{}{
		"result":      conc.Result,
		"description": conc.Description,
	}).Error
}

func (r *Repository) SetStatusFormed(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "черновик" {
		return fmt.Errorf("only draft can be formed")
	}
	var items []ds.ConcentrationItem
	if err := r.db.Where("concentration_id = ?", id).Find(&items).Error; err != nil {
		return err
	}
	result := r.calculateResult(items)
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":     "сформирован",
		"updated_at": time.Now(),
		"result":     result,
	}).Error
}

func (r *Repository) SetStatusFinished(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "сформирован" {
		return fmt.Errorf("only formed can be finished")
	}
	now := time.Now()
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":       "завершён",
		"updated_at":   now,
		"finished_at":  now,
		"moderator_id": 1, // TODO: брать из контекста модератора
	}).Error
}

func (r *Repository) SetStatusRejected(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "сформирован" {
		return fmt.Errorf("only formed can be rejected")
	}
	now := time.Now()
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":       "отклонён",
		"updated_at":   now,
		"finished_at":  now,
		"moderator_id": 1, // TODO: брать из контекста модератора
	}).Error
}

func (r *Repository) DeleteConcentration(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "черновик" && conc.Status != "сформирован" {
		return fmt.Errorf("only draft or formed can be deleted")
	}
	return r.db.Model(&conc).Update("status", "удалён").Error
}

func (r *Repository) calculateResult(items []ds.ConcentrationItem) string {
	// Здесь будет логика расчёта
	return "[H⁺] = 0.045 моль/л, pH = 1.35"
}