package usecase

import (
	"fmt"
	"music_catalog/internal/config"
	"music_catalog/internal/helper"
	"music_catalog/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type MembershipsRepo interface {
	CreateUsers(users models.User) error
	GetUsers(email, username string, id int) (users *models.User, err error)
}

type memberships struct {
	cfg  *config.Config
	repo MembershipsRepo
}

func NewMembershipUc(cfg *config.Config, repo MembershipsRepo) *memberships {
	return &memberships{
		cfg:  cfg,
		repo: repo,
	}
}

func (u *memberships) SignIn(req models.SignReq) error {
	_, err := u.repo.GetUsers(req.Email, req.Username, 0)

	if err == nil {
		return fmt.Errorf(helper.ErrExistUser)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(pass),
	}
	return u.repo.CreateUsers(user)
}
