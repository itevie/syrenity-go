package routes

import (
	"fmt"
	"net/http"
	"os"
	"syrenity/server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var CDNPath = "/home/isabella/Documents/projects/go/syrenity/files"

func RegisterCDNRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	router.GET("/files/:file_id", func(c *gin.Context) {
		file_id := c.Param("file_id")

		var file models.File
		err := db.QueryRowx("SELECT * FROM files WHERE id = $1;", file_id).StructScan(&file)

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, models.ErrorMessage{
				Message: "Failed to find file",
			})
			return
		}

		var path = fmt.Sprintf("%s/%s/%s-%s", CDNPath, file.CreatedAt.Format("02-01-2006"), file.ID, file.FileName)

		if _, err := os.Stat(path); err == nil {
			c.File(path)
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "The file exists, but it does not exist on disk",
			})
			return
		}
	})
}
