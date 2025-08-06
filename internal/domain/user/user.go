package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	SteamID   string    `db:"steam_id"`
	Username  string    `db:"username"`
	Cash      float64   `db:"cash"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	TradeUrl  string    `db:"trade_url"`
	AvatarURL string    `db:"avatar_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
