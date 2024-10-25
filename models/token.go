package models

import "time"

type Token struct {
	Token      string    `db:"token"`
	Account    int       `db:"account"`
	CreatedAt  time.Time `db:"created_at"`
	Identifier string    `db:"identifier"`
}
