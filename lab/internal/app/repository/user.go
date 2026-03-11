package repository

import (
	"lab/internal/app/ds"
)

func (r *Repository) GetUserByID(id uint) (*ds.User, error) {
	var user ds.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *Repository) CreateUser(login, password string) error {
	user := ds.User{
		Login:    login,
		Password: password, // в реальности нужно хешировать
	}
	return r.db.Create(&user).Error
}

func (r *Repository) Authenticate(login, password string) (uint, error) {
	var user ds.User
	err := r.db.Where("login = ? AND password = ?", login, password).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}