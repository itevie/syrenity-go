package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"syrenity/server/models"
	"syrenity/server/socket"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterChannelRoutes(router *gin.RouterGroup, db *sqlx.DB, ws *socket.WebsocketServer) {
	router.GET("/channels/:channel_id", func(c *gin.Context) {
		id := c.Param("channel_id")

		var channel models.Channel
		err := db.QueryRowx("SELECT * FROM channels WHERE id = $1;", id).StructScan(&channel)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to fetch channel",
			})
			return
		}

		c.JSON(http.StatusOK, channel)
	})

	router.GET("/channels/:channel_id/messages", func(c *gin.Context) {
		channel_id, _ := strconv.Atoi(c.Param("channel_id"))

		messages, err := models.MessageQuery{
			Amount:    20,
			ChannelID: channel_id,
		}.Query(db)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, messages)
	})

	router.POST("/channels/:channel_id/messages", func(c *gin.Context) {
		channel_id, _ := strconv.Atoi(c.Param("channel_id"))

		var user models.User
		models.UserFromAuth(db, c.GetHeader("Authorization"), &user)

		var message_opt struct {
			Content string
		}

		if err := c.ShouldBindBodyWithJSON(&message_opt); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: err.Error(),
			})
			return
		}

		message, err := CreateMessage(db, ws, MessageCreationOptions{
			Content:   message_opt.Content,
			AuthorId:  user.ID,
			ChannelID: channel_id,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to create message",
			})
			return
		}

		c.JSON(http.StatusOK, message)
	})
}

type MessageCreationOptions struct {
	ChannelID int
	AuthorId  int
	Content   string
}

func CreateMessage(db *sqlx.DB, ws *socket.WebsocketServer, options MessageCreationOptions) (*models.Message, error) {
	var message models.Message
	err := db.QueryRowx(
		"INSERT INTO messages (channel_id, author_id, content, created_at) VALUES ($1, $2, $3, $4) RETURNING *;",
		options.ChannelID, options.AuthorId, options.Content, time.Now(),
	).StructScan(&message)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	jsonText, _ := json.Marshal(message)

	ws.HandleMessage(socket.Message{
		Type:    socket.MTDMessageCreate,
		Payload: string(jsonText),
	})

	return &message, nil
}
