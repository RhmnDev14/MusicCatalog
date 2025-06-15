package models

import "github.com/golang-jwt/jwt/v5"

type (
	Claim struct {
		jwt.RegisteredClaims
		UserId int `json:"user_id"`
	}
	AuthResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		UserId       int    `json:"user_id"`
	}
)
