package internal

type Skins struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Weapon      struct {
		ID       string `json:"id"`
		WeaponID int    `json:"weapon_id"`
		Name     string `json:"name"`
	} `json:"weapon"`
	Category struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Pattern struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"pattern"`
	MinFloat float64 `json:"min_float"`
	MaxFloat float64 `json:"max_float"`
	Rarity   struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"rarity"`
	Stattrak   bool   `json:"stattrak"`
	Souvenir   bool   `json:"souvenir,omitempty"`
	PaintIndex string `json:"paint_index"`
	Wears      []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"wears,omitempty"`
	Collections []interface{} `json:"collections,omitempty"`
	Crates      []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"crates"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	LegacyModel  bool   `json:"legacy_model"`
	Image        string `json:"image"`
	Phase        string `json:"phase,omitempty"`
	SpecialNotes []struct {
		Source string `json:"source"`
		Text   string `json:"text"`
	} `json:"special_notes,omitempty"`
}

type InventoryResponse struct {
	Assets              []Asset       `json:"assets"`
	TotalInventoryCount int           `json:"total_inventory_count"`
	Success             int           `json:"success"`
	Descriptions        []Description `json:"descriptions"`
}

type Asset struct {
	Appid      int    `json:"appid"`
	Contextid  string `json:"contextid"`
	Assetid    string `json:"assetid"`
	Classid    string `json:"classid"`
	Instanceid string `json:"instanceid"`
	Amount     string `json:"amount"`
}

type Description struct {
	Appid           int    `json:"appid"`
	Classid         string `json:"classid"`
	Instanceid      string `json:"instanceid"`
	Currency        int    `json:"currency"`
	BackgroundColor string `json:"background_color"`
	IconURL         string `json:"icon_url"`
	Descriptions    []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Name  string `json:"name"`
		Color string `json:"color,omitempty"`
	} `json:"descriptions"`
	Tradable int `json:"tradable"`
	Actions  []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"actions"`
	Name           string `json:"name"`
	NameColor      string `json:"name_color"`
	Type           string `json:"type"`
	MarketName     string `json:"market_name"`
	MarketHashName string `json:"market_hash_name"`
	MarketActions  []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"market_actions"`
	Commodity                   int    `json:"commodity"`
	MarketTradableRestriction   int    `json:"market_tradable_restriction"`
	MarketMarketableRestriction int    `json:"market_marketable_restriction"`
	Marketable                  int    `json:"marketable"`
	Tags                        []Tags `json:"tags"`
	Sealed                      int    `json:"sealed"`
}
type Tags struct {
	Category              string `json:"category"`
	InternalName          string `json:"internal_name"`
	LocalizedCategoryName string `json:"localized_category_name"`
	LocalizedTagName      string `json:"localized_tag_name"`
	Color                 string `json:"color,omitempty"`
}
