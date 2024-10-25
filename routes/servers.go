package routes

import (
	"fmt"
	"net/http"
	"syrenity/server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterServerRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	router.GET("/servers/:server_id", func(c *gin.Context) {
		id := c.Param("server_id")

		var server models.Server
		err := db.QueryRowx("SELECT * FROM guilds WHERE id = $1;", id).StructScan(&server)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to fetch server",
			})
			return
		}

		c.JSON(http.StatusOK, server)
	})

	router.GET("/servers/:server_id/channels", func(c *gin.Context) {
		id := c.Param("server_id")

		var channels []models.Channel
		err := db.Select(&channels, "SELECT * FROM channels WHERE guild_id = $1;", id)

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to fetch server channels",
			})
			return
		}

		fmt.Println(channels)

		c.JSON(http.StatusOK, channels)
	})
}
