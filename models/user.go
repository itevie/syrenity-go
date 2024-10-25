package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID              int       `db:"id" json:"id"`
	Username        string    `db:"username" json:"username"`
	Password        string    `db:"password" json:"-"`
	Avatar          *string   `db:"avatar" json:"avatar"`
	TwoFactorSecret *string   `db:"2fa_secret" json:"-"`
	IsBot           bool      `db:"is_bot" json:"is_bot"`
	AboutMe         *string   `db:"about_me" json:"about_me"`
	Discriminator   string    `db:"discriminator" json:"discriminator"`
	Email           *string   `db:"email" json:"email"`
	EmailVerified   bool      `db:"email_verified" json:"email_verified"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func UserFromAuth(db *sqlx.DB, t string, w *User) error {
	var token Token
	err := db.QueryRowx("SELECT * FROM tokens WHERE token = $1;", strings.Replace(t, "Token ", "", 1)).StructScan(&token)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err2 := db.QueryRowx("SELECT * FROM users WHERE id = $1;", token.Account).StructScan(w)

	if err2 != nil {
		fmt.Println(err2.Error())
		return err
	}

	return nil
}
