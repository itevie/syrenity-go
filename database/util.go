package database

import (
	"syrenity/server/models"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	Inner *sqlx.DB
}

func (db *Database) GetUser(id int, user *models.User) error {
	return db.Inner.QueryRowx("SELECT * FROM users WHERE id = $1;", id).StructScan(&user)
}
