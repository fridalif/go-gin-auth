package handlers

import (
	"errors"
	"hitenok/pkg/config"
	"hitenok/pkg/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRequest struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}
type AuthHandlerI interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type MailAuthHandler struct {
	authenticationService services.PasswordAuthenticationServiceI
	otpService            services.OTPServiceI
	jwtService            services.JWTServiceI
	appConfig             *config.AppConfig
}

func NewMailAuthHandler(authenticationService services.PasswordAuthenticationServiceI, otpService services.OTPServiceI, jwtService services.JWTServiceI, appConfig *config.AppConfig) AuthHandlerI {
	return &MailAuthHandler{
		authenticationService: authenticationService,
		otpService:            otpService,
		appConfig:             appConfig,
		jwtService:            jwtService,
	}
}

func (mailAuthHandler *MailAuthHandler) SignIn(c *gin.Context) {
	var userRequest UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}

	user, err := mailAuthHandler.authenticationService.Authenticate(userRequest.Email, userRequest.Password)
	if errors.Is(err.ErrorBase, gorm.ErrRecordNotFound) || err.ErrorBase.Error() == "wrong credentials" || err.ErrorBase.Error() == "user is not active" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
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
		log.Printf("mailAuthHandler.SignIn.%s: %v", err.Module, err.ErrorBase)
		return
	}
	accessToken, err := mailAuthHandler.jwtService.GenerateToken(user, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("mailAuthHandler.SignIn.%s: %v", err.Module, err.ErrorBase)
		return
	}
	refreshToken, err := mailAuthHandler.jwtService.GenerateToken(user, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("mailAuthHandler.SignIn.%s: %v", err.Module, err.ErrorBase)
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
}

func (mailAuthHandler *MailAuthHandler) SignUp(c *gin.Context) {
	var userRequest UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong credentials",
		})
		return
	}

	user, err := mailAuthHandler.authenticationService.Register(userRequest.Email, userRequest.Fullname, userRequest.Password)
	if err.ErrorBase.Error() == "user already exists" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "User already exists",
		})
		return
	}
	if err.ErrorBase.Error() == "invalid credentials" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Invalid credentials",
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("mailAuthHandler.SignUp.%s: %v", err.Module, err.ErrorBase)
		return
	}
	err = mailAuthHandler.otpService.GenerateOTP(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("mailAuthHandler.SignUp.%s: %v", err.Module, err.ErrorBase)
		return
	}
	go mailAuthHandler.otpService.SendOTP(*user)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body": gin.H{
			"user_id": user.ID,
		},
		"error": nil,
	})

}

func (mailAuthHandler *MailAuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	mail := router.Group("/mail")
	mail.POST("/sign-in", mailAuthHandler.SignIn)
	mail.POST("/sign-up", mailAuthHandler.SignUp)
}
