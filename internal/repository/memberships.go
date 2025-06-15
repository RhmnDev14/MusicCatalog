package repository

import (
	"music_catalog/internal/models"

	"gorm.io/gorm"
)

type memberships struct {
	db *gorm.DB
}

func NewMembershipRepo(db *gorm.DB) *memberships {
	return &memberships{
		db: db,
	}
}

func (r *memberships) CreateUsers(users models.User) error {
	return r.db.Create(&users).Error
}

func (r *memberships) GetUsers(email, username string, id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("email = ? OR username = ? OR id = ?", email, username, id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
