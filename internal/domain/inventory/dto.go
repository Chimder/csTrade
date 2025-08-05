package inventory

type InventoryResponse struct {
	Success             bool          `json:"success"`
	TotalInventoryCount int           `json:"total_inventory_count"`
	Assets              []Asset       `json:"assets"`
	Descriptions        []Description `json:"descriptions"`
}

type Asset struct {
	AppID      int64  `json:"appid"`
	ContextID  string `json:"contextid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
	Amount     string `json:"amount"`
}

type Description struct {
	AppID                     int64  `json:"appid"`
	ClassID                   string `json:"classid"`
	InstanceID                string `json:"instanceid"`
	MarketHashName            string `json:"market_hash_name"`
	Name                      string `json:"name"`
	Tradable                  int    `json:"tradable"`
	Marketable                int    `json:"marketable"`
	MarketTradableRestriction int    `json:"market_tradable_restriction"`
	IconURL                   string `json:"icon_url"`
	IconURLLarge              string `json:"icon_url_large,omitempty"`
	Type                      string `json:"type"`
	MarketName                string `json:"market_name"`
	NameColor                 string `json:"name_color,omitempty"`
	BackgroundColor           string `json:"background_color,omitempty"`
	Rarity                    string `json:"rarity,omitempty"`
}
