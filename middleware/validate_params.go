package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"syrenity/server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func ValidateParams(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		authUserError := models.UserFromAuth(db, c.GetHeader("Authorization"), &user)

		userParam := c.Param("user_id")

		if userParam != "" {
			var _user models.User
			err := db.QueryRowx("SELECT * FROM users WHERE id = $1;", userParam).StructScan(&_user)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Invalid user param",
				})
				c.Abort()
				return
			}
		}

		serverParam := c.Param("server_id")

		if serverParam != "" {
			var server models.Server
			err := db.QueryRowx("SELECT * FROM guilds WHERE id = $1;", serverParam).StructScan(&server)

			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Invalid server param",
				})
				c.Abort()
				return
			}

			if authUserError != nil || !user.CanViewServer(strconv.Itoa(server.ID), db) {
				c.JSON(http.StatusForbidden, models.ErrorMessage{
					Message: "User cannot view this server",
				})
				c.Abort()
				return
			}
		}

		channelParam := c.Param("channel_id")

		if channelParam != "" {
			var channel models.Channel
			err := db.QueryRowx("SELECT * FROM channels WHERE id = $1;", channelParam).StructScan(&channel)

			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Invalid channel param",
				})
				c.Abort()
				return
			}
		}
	}
}
