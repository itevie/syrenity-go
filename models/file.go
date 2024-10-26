package models

import "time"

type File struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	FileName  string    `db:"file_name" json:"file_name"`
}
