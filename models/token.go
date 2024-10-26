package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Token      string    `db:"token"`
	Account    int       `db:"account"`
	CreatedAt  time.Time `db:"created_at"`
	Identifier string    `db:"identifier"`
}

func GenerateToken() string {
	return uuid.New().String()
}
