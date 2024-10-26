package models

type Server struct {
	ID           int     `db:"id" json:"id"`
	Name         string  `db:"name" json:"name"`
	OwnerId      int     `db:"owner_id" json:"owner_id"`
	Description  *string `db:"description" json:"description"`
	Avatar       *string `db:"avatar" json:"avatar"`
	ChannelOrder []uint8 `db:"channel_order" json:"channel_order"`
}
