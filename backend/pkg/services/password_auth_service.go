package services

import (
	"errors"
	"fmt"
	"hitenok/pkg/config"
	"hitenok/pkg/domain"
	"hitenok/pkg/repository"
	"hitenok/pkg/security"

	"gorm.io/gorm"
)

type PasswordAuthenticationServiceI interface {
	Authenticate(credentials, password string) (*domain.User, *domain.MyError)
	Register(credentials, fullname, password string) (*domain.User, *domain.MyError)
}

type mailAuthenticationService struct {
	repo      repository.UserRepositoryI
	appConfig *config.AppConfig
}

func NewMailAuthenticationService(repo repository.UserRepositoryI, appConfig *config.AppConfig) PasswordAuthenticationServiceI {
	return &mailAuthenticationService{
		repo:      repo,
		appConfig: appConfig,
	}
}

func (mailAuthenticationService *mailAuthenticationService) Authenticate(email, password string) (*domain.User, *domain.MyError) {
	user, err := mailAuthenticationService.repo.FindUserByEmail(email)
	if err != nil {
		err.Module = "mailAuthenticationService.Authenticate" + err.Module
		return user, err
	}
	if !user.IsActive {
		return user, domain.NewError(fmt.Errorf("user is not active"), "mailAuthenticationService.Authenticate")
	}
	if user.Password == security.HashPassword(password, mailAuthenticationService.appConfig.SecretKey) {
		return user, nil
	}
	return user, domain.NewError(fmt.Errorf("wrong credentials"), "mailAuthenticationService.Authenticate")
}

func (mailAuthenticationService *mailAuthenticationService) Register(email, fullname, password string) (*domain.User, *domain.MyError) {
	if (email == "") || (fullname == "") || (password == "") {
		return &domain.User{}, domain.NewError(fmt.Errorf("invalid credentials"), "mailAuthenticationService.Register")
	}
	_, err := mailAuthenticationService.repo.FindUserByEmail(email)
	if err == nil {
		return &domain.User{}, domain.NewError(fmt.Errorf("user already exists"), "mailAuthenticationService.Register")
	}
	if !errors.Is(err.ErrorBase, gorm.ErrRecordNotFound) {
		return &domain.User{}, err
	}
	user := &domain.User{
		Email:    email,
		Fullname: fullname,
		Password: security.HashPassword(password, mailAuthenticationService.appConfig.SecretKey),
	}
	err = mailAuthenticationService.repo.SaveUser(user)
	if err != nil {
		err.Module = "mailAuthenticationService.Register" + err.Module
		return user, err
	}
	return user, nil
}
