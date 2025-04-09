package services

import (
	"hitenok/pkg/domain"
	"hitenok/pkg/repository"
	"math/rand"
	"time"
)

const hashCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type HashServiceI interface {
	GenerateHash(user *domain.User) *domain.MyError
	ValidateHash(user *domain.User, hash string) (bool, *domain.MyError)
	ClearHash(user *domain.User) *domain.MyError
}

type HashService struct {
	userRepo repository.UserRepositoryI
}

func NewHashService(userRepo repository.UserRepositoryI) HashServiceI {
	return &HashService{
		userRepo: userRepo,
	}
}

func (hashService *HashService) GenerateHash(user *domain.User) *domain.MyError {
	newHash := make([]byte, 16)
	for i := range newHash {
		newHash[i] = hashCharset[rand.Intn(len(hashCharset))]
	}
	user.ResetHash = string(newHash)
	user.ResetHashSpawnedAt = time.Now()
	user.HashAttempts = 3
	err := hashService.userRepo.SaveUser(user)
	if err != nil {
		err.Module = "HashService.GenerateHash." + err.Module
		return err
	}
	return nil
}

func (hashService *HashService) ValidateHash(user *domain.User, hash string) (bool, *domain.MyError) {
	if user.ResetHash == "" {
		return false, nil
	}
	if user.ResetHashSpawnedAt.Add(5 * time.Minute).Before(time.Now()) {
		return false, nil
	}
	if user.HashAttempts <= 0 {
		return false, nil
	}
	if hash != user.ResetHash {
		user.HashAttempts -= 1
		err := hashService.userRepo.SaveUser(user)
		if err != nil {
			err.Module = "HashService.ValidateHash." + err.Module
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (hashService *HashService) ClearHash(user *domain.User) *domain.MyError {
	user.HashAttempts = 0
	user.ResetHash = ""
	user.ResetHashSpawnedAt = time.Time{}
	err := hashService.userRepo.SaveUser(user)
	if err != nil {
		err.Module = "HashService.ClearHash." + err.Module
		return err
	}
	return nil
}
