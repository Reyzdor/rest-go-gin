package middleware

import (
	"Application/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")

		if err != nil {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token != "" {
			claims, err := auth.ValidateAccessToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("user_email", claims.Email)
				c.Set("is_authenticated", true)
			} else {
				c.Set("is_authenticated", false)
			}
		} else {
			c.Set("is_authenticated", false)
		}

		c.Next()
	}
}
