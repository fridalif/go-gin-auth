package services

import (
	"hitenok/pkg/domain"
	"hitenok/pkg/repository"
)

type UserServiceI interface {
	GetUser(id uint) (*domain.User, *domain.MyError)
	GetUserByEmail(email string) (*domain.User, *domain.MyError)
	UpdateUser(user *domain.User) *domain.MyError
}

type userService struct {
	repo repository.UserRepositoryI
}

func NewUserService(repo repository.UserRepositoryI) UserServiceI {
	return &userService{
		repo: repo,
	}
}

func (userService *userService) GetUser(id uint) (*domain.User, *domain.MyError) {
	user, err := userService.repo.FindUserById(id)
	if err != nil {
		err.Module = "userService.GetUser" + err.Module
		return user, err
	}
	return user, nil
}

func (userService *userService) GetUserByEmail(email string) (*domain.User, *domain.MyError) {
	user, err := userService.repo.FindUserByEmail(email)
	if err != nil {
		err.Module = "userService.GetUserByEmail" + err.Module
		return user, err
	}
	return user, nil
}

func (userService *userService) UpdateUser(user *domain.User) *domain.MyError {
	err := userService.repo.SaveUser(user)
	if err != nil {
		err.Module = "userService.UpdateUser" + err.Module
		return err
	}
	return nil
}
