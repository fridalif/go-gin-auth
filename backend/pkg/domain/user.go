package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string    `json:"email" gorm:"not null"`
	Password           string    `json:"password" gorm:"not null"`
	Fullname           string    `json:"fullname" gorm:"not null"`
	IsSuperuser        bool      `json:"isSuperuser" gorm:"default:false"`
	OTP                string    `json:"otp"`
	OTPAttempts        int       `json:"otpAttempts" gorm:"default:0"`
	OTPSpawnedAt       time.Time `json:"otpSpawnedAt"`
	ResetHash          string    `json:"resetHash"`
	HashAttempts       int       `json:"hashAttempts" gorm:"default:0"`
	ResetHashSpawnedAt time.Time `json:"resetHashSpawnedAt"`
	IsActive           bool      `json:"isActive" gorm:"default:false"`
	JWTVersion         uint      `json:"jwtVersion" gorm:"default:0"`
}
