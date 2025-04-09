package handlers

import (
	"errors"
	"hitenok/pkg/config"
	"hitenok/pkg/domain"
	"hitenok/pkg/security"
	"hitenok/pkg/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivateRequest struct {
	UserId      uint   `json:"user_id"`
	Email       string `json:"email"`
	OTP         string `json:"one_time_password"`
	ResetHash   string `json:"reset_hash"`
	NewPassword string `json:"new_password"`
}

type ActivateHandlerI interface {
	Activate(c *gin.Context)
	Resend(c *gin.Context)
	ResetPassword(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type ActivateHandler struct {
	otpService  services.OTPServiceI
	userService services.UserServiceI
	hashService services.HashServiceI
	jwtService  services.JWTServiceI
	appConfig   *config.AppConfig
}

func NewActivateHandler(otpService services.OTPServiceI, hashService services.HashServiceI, jwtService services.JWTServiceI, userService services.UserServiceI, appConfig *config.AppConfig) ActivateHandlerI {
	return &ActivateHandler{
		otpService:  otpService,
		hashService: hashService,
		userService: userService,
		jwtService:  jwtService,
		appConfig:   appConfig,
	}
}

func (activateHandler *ActivateHandler) Activate(c *gin.Context) {
	var activateRequest ActivateRequest
	if err := c.ShouldBindJSON(&activateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	if activateRequest.UserId == 0 || activateRequest.OTP == "" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	user, err := activateHandler.userService.GetUser(activateRequest.UserId)
	if errors.Is(err.ErrorBase, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
		return
	}
	valid, err := activateHandler.otpService.VerifyOTP(*user, activateRequest.OTP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
		return
	}
	if !valid {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	if !user.IsActive {
		user.IsActive = true
		err := activateHandler.userService.UpdateUser(user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"body":   gin.H{},
				"error":  "Internal server error",
			})
			log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
			return
		}
		accessToken, err := activateHandler.jwtService.GenerateToken(user, true)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"body":   gin.H{},
				"error":  "Internal server error",
			})
			log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
			return
		}
		refreshToken, err := activateHandler.jwtService.GenerateToken(user, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"body":   gin.H{},
				"error":  "Internal server error",
			})
			log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"body": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
			"error": nil,
		})
		go func() {
			err := activateHandler.otpService.ClearOTP(user)
			if err != nil {
				log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
			}
		}()
	}
	err = activateHandler.hashService.GenerateHash(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body": gin.H{
			"reset_hash": user.ResetHash,
		},
		"error": nil,
	})
	go func() {
		err := activateHandler.otpService.ClearOTP(user)
		if err != nil {
			log.Printf("activateHandler.Activate.%s: %v", err.Module, err.ErrorBase)
		}
	}()
}

func (activateHandler *ActivateHandler) Resend(c *gin.Context) {
	var activateRequest ActivateRequest
	if err := c.ShouldBindJSON(&activateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	var user *domain.User
	var err *domain.MyError
	if activateRequest.Email != "" {
		user, err = activateHandler.userService.GetUserByEmail(activateRequest.Email)
	} else {
		user, err = activateHandler.userService.GetUser(activateRequest.UserId)
	}
	if errors.Is(err.ErrorBase, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.Resend.%s: %v", err.Module, err.ErrorBase)
		return
	}

	err = activateHandler.otpService.GenerateOTP(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.Resend.%s: %v", err.Module, err.ErrorBase)
		return
	}
	go activateHandler.otpService.SendOTP(*user)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body": gin.H{
			"user_id": user.ID,
		},
		"error": nil,
	})
}

func (activateHandler *ActivateHandler) ResetPassword(c *gin.Context) {
	var activateRequest ActivateRequest
	if err := c.ShouldBindJSON(&activateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}

	user, err := activateHandler.userService.GetUser(activateRequest.UserId)
	if errors.Is(err.ErrorBase, gorm.ErrRecordNotFound) || activateRequest.ResetHash == "" || activateRequest.NewPassword == "" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.ResetPassword.%s: %v", err.Module, err.ErrorBase)
		return
	}
	user.Password = security.HashPassword(activateRequest.NewPassword, activateHandler.appConfig.SecretKey)
	user.JWTVersion += 1
	err = activateHandler.userService.UpdateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("activateHandler.ResetPassword.%s: %v", err.Module, err.ErrorBase)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body":   gin.H{},
		"error":  nil,
	})
	go func() {
		err := activateHandler.hashService.ClearHash(user)
		if err != nil {
			log.Printf("activateHandler.ResetPassword.%s: %v", err.Module, err.ErrorBase)
		}
	}()
}

func (activateHandler *ActivateHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/activate", activateHandler.Activate)
	router.POST("/resend", activateHandler.Resend)
	router.POST("/reset-password", activateHandler.ResetPassword)
}
