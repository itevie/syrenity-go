package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	ID         int       `db:"id" json:"id"`
	ChannelID  int       `db:"channel_id" json:"channel_id"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	AuthorId   int       `db:"author_id" json:"author_id"`
	IsPinned   bool      `db:"is_pinned" json:"is_pinned"`
	IsEdited   bool      `db:"is_edited" json:"is_edited"`
	IsSystem   bool      `db:"is_system" json:"is_system"`
	SystemType *string   `db:"sys_type" json:"sys_type"`
}

type MessageQuery struct {
	Amount    int
	ChannelID int
}

func (query MessageQuery) Query(db *sqlx.DB) ([]Message, error) {
	var baseQuery = fmt.Sprintf("SELECT * FROM messages WHERE channel_id = %d LIMIT %d", query.ChannelID, query.Amount)

	var messages []Message
	err := db.Select(&messages, baseQuery)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New(err.Error())
	}

	return messages, nil
}
