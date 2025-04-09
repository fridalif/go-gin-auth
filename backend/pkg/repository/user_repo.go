package repository

import (
	"hitenok/pkg/config"
	"hitenok/pkg/domain"

	"gorm.io/gorm"
)

type UserRepositoryI interface {
	FindUserById(id uint) (*domain.User, *domain.MyError)
	FindUserByEmail(email string) (*domain.User, *domain.MyError)
	SaveUser(user *domain.User) *domain.MyError
}

type userRepository struct {
	DB        *gorm.DB
	AppConfig *config.AppConfig
}

func NewUserRepository(db *gorm.DB, appConfig *config.AppConfig) UserRepositoryI {
	return &userRepository{
		DB:        db,
		AppConfig: appConfig,
	}
}

func (userRepo *userRepository) FindUserById(id uint) (*domain.User, *domain.MyError) {
	var user domain.User
	err := userRepo.DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return &user, domain.NewError(err, "userRepository.FindUserById")
	}
	return &user, nil
}

func (userRepository *userRepository) FindUserByEmail(email string) (*domain.User, *domain.MyError) {
	var user domain.User
	err := userRepository.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return &user, domain.NewError(err, "userRepository.FindUserByEmail")
	}
	return &user, nil
}

func (baseRepo *userRepository) SaveUser(user *domain.User) *domain.MyError {
	err := baseRepo.DB.Save(&user).Error
	if err != nil {
		return domain.NewError(err, "userRepository.SaveUser")
	}
	return nil
}
