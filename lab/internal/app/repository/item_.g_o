package repository

import (
	"fmt"
	"lab/internal/app/ds"
)

func (r *Repository) AddItemToDraft(electrolyteID uint, volume int, comment string) error {
	draft, err := r.GetOrCreateDraft()
	if err != nil {
		return err
	}
	// Проверяем, нет ли уже такого электролита
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

func (r *Repository) UpdateItem(itemID uint, volume int, comment string) error {
	return r.db.Model(&ds.ConcentrationItem{}).Where("id = ?", itemID).Updates(map[string]interface{}{
		"volume":  volume,
		"comment": comment,
	}).Error
}

func (r *Repository) DeleteItem(itemID uint) error {
	return r.db.Delete(&ds.ConcentrationItem{}, itemID).Error
}