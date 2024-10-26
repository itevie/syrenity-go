package models

type Member struct {
	GuildId  int     `db:"guild_id" json:"guild_id"`
	UserId   int     `db:"user_id" json:"user_id"`
	Nickname *string `db:"nickname" json:"nickname"`
}
