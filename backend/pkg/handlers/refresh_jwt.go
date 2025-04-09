package handlers

import (
	"hitenok/pkg/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshJWTHandler(c *gin.Context, jwtService services.JWTServiceI) {
	token := c.Request.Header.Get("Authorization")
	user, err := jwtService.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"body":   gin.H{},
			"error":  "Unauthorized",
		})
		return
	}
	accessToken, err := jwtService.GenerateToken(user, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("RefreshJWTHandler.jwtService.GenerateToken.%s: %v", err.Module, err.ErrorBase)
		return
	}
	refreshToken, err := jwtService.GenerateToken(user, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"body":   gin.H{},
			"error":  "Internal server error",
		})
		log.Printf("RefreshJWTHandler.jwtService.GenerateToken.%s: %v", err.Module, err.ErrorBase)
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
