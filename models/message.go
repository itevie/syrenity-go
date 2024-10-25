package models

import (
	"time"
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
