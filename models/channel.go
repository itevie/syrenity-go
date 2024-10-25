package models

type Channel struct {
	Id      int     `db:"id" json:"id"`
	Type    string  `db:"type" json:"text"`
	GuildId *int    `db:"guild_id" json:"guild_id"`
	Name    *string `db:"name" json:"name"`
	Topic   *string `db:"topic" json:"topic"`
	Nsfw    bool    `db:"is_nsfw" json:"is_nsfw"`
}
