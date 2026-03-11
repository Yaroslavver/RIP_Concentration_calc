package repository

import (
	"fmt"
	"lab/internal/app/ds"
	"lab/internal/app/usercontext"
	"time"
)

// GetOrCreateDraft возвращает черновик текущего пользователя, при отсутствии создаёт новый
func (r *Repository) GetOrCreateDraft() (*ds.Concentration, error) {
	userID := usercontext.GetCurrentUserID()
	var conc ds.Concentration
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&conc).Error
	if err == nil {
		return &conc, nil
	}
	// создаём новый
	conc = ds.Concentration{
		Status:    "черновик",
		CreatedAt: time.Now(),
		CreatorID: userID,
	}
	err = r.db.Create(&conc).Error
	if err != nil {
		return nil, err
	}
	return &conc, nil
}

// GetCartInfo возвращает ID черновика и количество позиций
func (r *Repository) GetCartInfo() (draftID uint, count int64, err error) {
	userID := usercontext.GetCurrentUserID()
	var conc ds.Concentration
	err = r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&conc).Error
	if err != nil {
		return 0, 0, nil // нет черновика
	}
	var cnt int64
	r.db.Model(&ds.ConcentrationItem{}).Where("concentration_id = ?", conc.ID).Count(&cnt)
	return conc.ID, cnt, nil
}

// GetConcentrationByID возвращает заявку вместе с элементами и данными электролитов
func (r *Repository) GetConcentrationByID(id uint) (*ds.Concentration, []ds.ConcentrationItem, error) {
	var conc ds.Concentration
	err := r.db.Preload("Creator").Preload("Moderator").First(&conc, id).Error
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

// GetUserConcentrations возвращает список заявок пользователя (кроме черновика и удалённых) с фильтрацией
func (r *Repository) GetUserConcentrations(status string, fromDate, toDate *time.Time) ([]ds.Concentration, error) {
	userID := usercontext.GetCurrentUserID()
	query := r.db.Where("creator_id = ? AND status NOT IN (?, ?)", userID, "черновик", "удалён")
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

// UpdateConcentration обновляет поля заявки (не статус)
func (r *Repository) UpdateConcentration(conc *ds.Concentration) error {
	return r.db.Model(conc).Updates(map[string]interface{}{
		"result": conc.Result,
	}).Error
}

// SetStatusFormed устанавливает статус "сформирован", вычисляет результат
func (r *Repository) SetStatusFormed(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "черновик" {
		return fmt.Errorf("only draft can be formed")
	}
	// Получаем элементы для расчёта
	var items []ds.ConcentrationItem
	if err := r.db.Where("concentration_id = ?", id).Find(&items).Error; err != nil {
		return err
	}
	result := r.calculateResult(items)
	// Обновляем
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":     "сформирован",
		"updated_at": time.Now(),
		"result":     result,
	}).Error
}

// SetStatusFinished завершает заявку (модератор)
func (r *Repository) SetStatusFinished(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "сформирован" {
		return fmt.Errorf("only formed can be finished")
	}
	if !usercontext.GetCurrentUserIsModerator() {
		return fmt.Errorf("only moderator can finish")
	}
	now := time.Now()
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":       "завершён",
		"updated_at":   now,
		"finished_at":  now,
		"moderator_id": usercontext.GetCurrentUserID(),
	}).Error
}

// SetStatusRejected отклоняет заявку (модератор)
func (r *Repository) SetStatusRejected(id uint) error {
	var conc ds.Concentration
	if err := r.db.First(&conc, id).Error; err != nil {
		return err
	}
	if conc.Status != "сформирован" {
		return fmt.Errorf("only formed can be rejected")
	}
	if !usercontext.GetCurrentUserIsModerator() {
		return fmt.Errorf("only moderator can reject")
	}
	now := time.Now()
	return r.db.Model(&conc).Updates(map[string]interface{}{
		"status":       "отклонён",
		"updated_at":   now,
		"finished_at":  now,
		"moderator_id": usercontext.GetCurrentUserID(),
	}).Error
}

// DeleteConcentration логическое удаление (статус "удалён")
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

// calculateResult - формула расчета итоговой концентрации (заглушка, но можно реализовать реальную)
func (r *Repository) calculateResult(items []ds.ConcentrationItem) string {
	// Здесь должна быть логика вычисления pH по объёмам и концентрациям
	// Пока возвращаем пример
	return "[H⁺] = 0.045 моль/л, pH = 1.35"
}