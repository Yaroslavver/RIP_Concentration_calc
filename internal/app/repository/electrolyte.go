package repository

import (
	"database/sql"
	"errors"
	"go_project2/internal/app/ds"
)

// GetAllElectrolytes возвращает все не удалённые растворы
func (r *Repository) GetAllElectrolytes() ([]ds.Electrolyte, error) {
	var electrolytes []ds.Electrolyte
	err := r.db.Where("is_deleted = ?", false).Find(&electrolytes).Error
	return electrolytes, err
}

// SearchElectrolytesByName ищет по названию (частичное совпадение)
func (r *Repository) SearchElectrolytesByName(name string) ([]ds.Electrolyte, error) {
	var electrolytes []ds.Electrolyte
	err := r.db.Where("name ILIKE ? AND is_deleted = ?", "%"+name+"%", false).Find(&electrolytes).Error
	return electrolytes, err
}

// GetElectrolyteByID использует "курсор" (raw SQL) для демонстрации
func (r *Repository) GetElectrolyteByID(id int) (*ds.Electrolyte, error) {
	query := `SELECT id, name, concentration, ions, ph, description, image, video 
	          FROM electrolytes WHERE id = $1 AND is_deleted = false`
	row := r.db.Raw(query, id).Row()
	var e ds.Electrolyte
	err := row.Scan(&e.ID, &e.Name, &e.Concentration, &e.Ions, &e.PH, &e.Description, &e.Image, &e.Video)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}