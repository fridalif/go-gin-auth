package handlers

import (
	"hitenok/pkg/domain"
	"hitenok/pkg/middlewares"
	"hitenok/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandlerI interface {
	UserInfo(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type UserHandler struct {
	userService services.UserServiceI
	jwtService  services.JWTServiceI
}

func NewUserHandler(userService services.UserServiceI, jwtService services.JWTServiceI) UserHandlerI {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (userHandler *UserHandler) UserInfo(c *gin.Context) {
	userInterface, exists := c.Get("user")
	user, ok := userInterface.(*domain.User)
	if !exists || !ok {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"body":   gin.H{},
			"error":  "Unauthorized",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body": gin.H{
			"user": user,
		},
		"error": nil,
	})
}

func (userHandler *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	user := router.Group("/users")
	user.Use(middlewares.CheckAuth(userHandler.jwtService))
	user.GET("/me", userHandler.UserInfo)
}
