package services

import (
	"fmt"
	"hitenok/pkg/config"
	"hitenok/pkg/domain"
	"hitenok/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceI interface {
	GenerateToken(user *domain.User, isAccess bool) (string, *domain.MyError)
	ValidateToken(token string) (*domain.User, *domain.MyError)
}

type JWTService struct {
	appConfig *config.AppConfig
	userRepo  repository.UserRepositoryI
}

func NewJWTService(appConfig *config.AppConfig, userRepo repository.UserRepositoryI) JWTServiceI {
	return &JWTService{
		appConfig: appConfig,
		userRepo:  userRepo,
	}
}

func (jwtService *JWTService) GenerateToken(user *domain.User, isAccess bool) (string, *domain.MyError) {
	expireTime := time.Now().Add(1 * time.Hour)
	if !isAccess {
		expireTime = time.Now().AddDate(0, 1, 0)
	}
	claims := domain.Claims{
		UserId:  user.ID,
		Version: user.JWTVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtService.appConfig.SecretKey))
	if err != nil {
		return "", domain.NewError(err, "JWTService.GenerateToken")
	}

	return tokenString, nil
}

func (jwtService *JWTService) ValidateToken(token string) (*domain.User, *domain.MyError) {
	tokenClaims := &domain.Claims{}
	_, err := jwt.ParseWithClaims(token, tokenClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtService.appConfig.SecretKey), nil
	})
	if err != nil {
		return nil, domain.NewError(err, "JWTService.ValidateToken")
	}
	user, customErr := jwtService.userRepo.FindUserById(tokenClaims.UserId)
	if customErr != nil {
		customErr.Module = "JWTService.ValidateToken" + customErr.Module
		return nil, customErr
	}
	if user.JWTVersion != tokenClaims.Version {
		return nil, domain.NewError(fmt.Errorf("invalid token"), "JWTService.ValidateToken")
	}
	return user, nil

}
