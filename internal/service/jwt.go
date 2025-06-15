package service

import (
	"fmt"
	"music_catalog/internal/config"
	"music_catalog/internal/helper"
	"music_catalog/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	cfgToken config.TokenConfig
}

func NewJwtService(cfgToken config.TokenConfig) *jwtService {
	return &jwtService{
		cfgToken: cfgToken,
	}
}

func (j *jwtService) CreateToken(user models.User) (*models.AuthResponse, error) {
	claims := models.Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfgToken.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfgToken.JwtExpiresTime)),
		},
		UserId: int(user.ID),
	}

	token := jwt.NewWithClaims(j.cfgToken.JwtSigningMethod, claims)

	ss, err := token.SignedString(j.cfgToken.JwtSignatureKey)
	if err != nil {
		return nil, fmt.Errorf("%s", helper.ErrCreateToken)
	}

	refreshToken, err := j.CreateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("%s", helper.ErrCreateTokenRefresh)
	}

	return &models.AuthResponse{
		Token:        ss,
		RefreshToken: refreshToken,
		UserId:       claims.UserId,
	}, nil
}

func (j *jwtService) CreateRefreshToken(user models.User) (string, error) {
	claims := models.Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfgToken.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfgToken.RefreshTokenExpiresTime)),
		},
		UserId: int(user.ID),
	}

	refreshToken := jwt.NewWithClaims(j.cfgToken.JwtSigningMethod, claims)

	return refreshToken.SignedString(j.cfgToken.JwtSignatureKey)
}
