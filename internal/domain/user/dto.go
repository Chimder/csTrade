package user

type UserCreateReq struct {
	SteamID   uint64 `db:"steam_id"`
	Username  string `db:"username"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	TradeUrl  string `db:"trade_url"`
	AvatarURL string `db:"avatar_url"`
}
