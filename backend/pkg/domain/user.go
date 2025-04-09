package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string    `json:"email" gorm:"not null"`
	Password           string    `json:"-" gorm:"not null"`
	Fullname           string    `json:"fullname" gorm:"not null"`
	IsSuperuser        bool      `json:"isSuperuser" gorm:"default:false"`
	OTP                string    `json:"-"`
	OTPAttempts        int       `json:"-" gorm:"default:0"`
	OTPSpawnedAt       time.Time `json:"-"`
	ResetHash          string    `json:"-"`
	HashAttempts       int       `json:"-" gorm:"default:0"`
	ResetHashSpawnedAt time.Time `json:"-"`
	IsActive           bool      `json:"isActive" gorm:"default:false"`
	JWTVersion         uint      `json:"jwtVersion" gorm:"default:0"`
}
