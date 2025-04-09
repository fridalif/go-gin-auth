package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId  uint `json:"user_id"`
	Version uint `json:"version"`
	jwt.RegisteredClaims
}
