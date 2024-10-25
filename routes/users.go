package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"syrenity/server/models"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	router.GET("/users/:user_id", func(c *gin.Context) {
		id := c.Param("user_id")

		var user models.User
		err := db.QueryRowx("SELECT * FROM users WHERE id = $1;", id).StructScan(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to fetch user",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.GET("/users/:user_id/servers", func(c *gin.Context) {
		id := c.Param("user_id")
		var auth models.User
		err := models.UserFromAuth(db, c.GetHeader("Authorization"), &auth)

		if err != nil || strconv.Itoa(auth.ID) != id {
			c.JSON(http.StatusForbidden, models.ErrorMessage{
				Message: "Expected user_id to be same as logged in user",
			})
			return
		}

		var servers []models.Server
		err = db.Select(&servers, `
			WITH guild_ids AS (
					SELECT guild_id
					FROM members
					WHERE user_id = $1
			)
			
			SELECT * 
					FROM guilds 
					WHERE (
							SELECT 1 
									FROM guild_ids 
									WHERE guild_id = guilds.id
											) = 1
		`, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to fetch user servers",
			})
			return
		}

		c.JSON(http.StatusOK, servers)
	})
}
