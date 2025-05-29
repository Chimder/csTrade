package bots

import (
	"time"
)

type SteamBot struct {
	AccountName    string    `json:"account_name"`
	Password       string    `json:"password"`
	SteamID        string    `json:"steam_id"`
	SharedSecret   string    `json:"shared_secret"`
	IdentitySecret string    `json:"identity_secret"`
	DeviceID       string    `json:"device_id"`
	AccessToken    string    `json:"access_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	RefreshToken   string    `json:"refresh_token"`
	RevocationCode string    `json:"revocation_code"`
}
