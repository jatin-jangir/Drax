package middleware

import (
	"alert-handler/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header must be Bearer <token>"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		c.Next()
	}
}
