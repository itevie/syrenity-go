package middleware

import (
	"net/http"
	"strings"

	"syrenity/server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RequireToken(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Check if it was present
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Missing Authorization field",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		// Check format
		if len(parts) != 2 || parts[0] != "Token" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Authorization field format",
			})
			c.Abort()
			return
		}

		// Get the by token
		var token models.Token
		err := db.QueryRowx("SELECT * FROM tokens WHERE token = $1;", parts[1]).StructScan(&token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid token provided",
			})
			c.Abort()
			return
		}
	}
}
