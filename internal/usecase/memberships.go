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

type jwtService interface {
	CreateToken(user models.User) (*models.AuthResponse, error)
}

type memberships struct {
	cfg  *config.Config
	repo MembershipsRepo
	jwt  jwtService
}

func NewMembershipUc(cfg *config.Config, repo MembershipsRepo, jwt jwtService) *memberships {
	return &memberships{
		cfg:  cfg,
		repo: repo,
		jwt:  jwt,
	}
}

func (u *memberships) SignIn(req models.SignReq) error {
	_, err := u.repo.GetUsers(req.Email, req.Username, 0)
	if err == nil {
		return fmt.Errorf("%s", helper.ErrExistUser)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s", helper.ErrGeneratePassword)
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(pass),
	}
	return u.repo.CreateUsers(user)
}

func (u *memberships) LogIn(req models.LogInReq) (*models.AuthResponse, error) {
	user, err := u.repo.GetUsers(req.Email, "", 0)
	if err != nil {
		return nil, fmt.Errorf("%s", helper.ErrGetUser)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("%s", helper.ErrComparePassword)
	}

	authResp, err := u.jwt.CreateToken(*user)
	if err != nil {
		return nil, err
	}

	return authResp, nil
}
