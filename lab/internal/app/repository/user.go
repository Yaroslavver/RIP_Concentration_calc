package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"lab/internal/app/ds"
)

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*ds.User, error) {
	var user ds.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *Repository) CreateUser(login, password string) error {
	hashed := hashPassword(password)
	user := ds.User{
		Login:    login,
		Password:    hashed,
		IsModerator: false,
	}
	return r.db.Create(&user).Error
}

func hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(hash[:])
}

func (r *Repository) CheckPassword(user *ds.User, password string) bool {
	return user.Password == hashPassword(password)
}

/*func (r *Repository) Authenticate(login, password string) (uint, error) {
	var user ds.User
	err := r.db.Where("login = ? AND password = ?", login, password).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}*/