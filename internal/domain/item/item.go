package item

import "time"

type Item struct {
	Id        string    `json:"id"`
	OwnerId   string    `json:"owner_id"`
	Price     int       `json:"price"`
	BotId     string    `json:"bot_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Float     float64   `json:"float"`
}
