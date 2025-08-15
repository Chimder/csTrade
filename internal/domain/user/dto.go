package user

type UserCreateReq struct {
	SteamID   string `json:"steam_id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	TradeUrl  string `json:"trade_url"`
	AvatarURL string `json:"avatar_url"`
}
