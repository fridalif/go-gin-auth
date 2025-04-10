package middlewares

import (
	"errors"
	"hitenok/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(jwtService services.JWTServiceI) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusUnauthorized,
				"body":   gin.H{},
				"error":  "Unauthorized",
			})
			return
		}
		user, err := jwtService.ValidateToken(token)
		if err != nil && errors.Is(err.ErrorBase, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusUnauthorized,
				"body":   gin.H{},
				"error":  "Token expired",
			})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": http.StatusUnauthorized,
				"body":   gin.H{},
				"error":  "Unauthorized",
			})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
