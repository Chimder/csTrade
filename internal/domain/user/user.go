package user

import (
	"time"
)

type UserDB struct {
	SteamID   uint64    `db:"steam_id"`
	Username  string    `db:"username"`
	Cash      float64   `db:"cash"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	TradeUrl  string    `db:"trade_url"`
	AvatarURL string    `db:"avatar_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
