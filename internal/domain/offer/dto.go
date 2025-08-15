package offer

type OfferCreateReq struct {
	SellerID   string  `json:"seller_id"`
	BotSteamID string  `json:"bot_steam_id"`
	Price      float64 `json:"price"`

	AssetID    string `json:"asset_id"`
	ClassID    string `json:"class_id"`
	InstanceID string `json:"instance_id"`

	Name                      string  `json:"name"`
	FullName                  string  `json:"full_name"`
	MarketTradableRestriction int     `json:"market_tradable_restriction"`
	IconURL                   string  `json:"icon_url"`
	NameColor                 string  `json:"name_color"`
	ActionLink                *string `json:"action_link"`
	TagType                   string  `json:"tag_type"`
	TagWeaponInternal         string  `json:"tag_weapon_internal"`
	TagWeaponName             string  `json:"tag_weapon_name"`
	TagQuality                string  `json:"tag_quality"`
	TagRarity                 string  `json:"tag_rarity"`
	TagRarityColor            string  `json:"tag_rarity_color"`
	TagExterior               string  `json:"tag_exterior"`
}
