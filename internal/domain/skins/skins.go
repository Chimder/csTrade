package skins

type Skins struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    struct {
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
	Stattrak    bool   `json:"stattrak"`
	Souvenir    bool   `json:"souvenir,omitempty"`
	PaintIndex  string `json:"paint_index"`
	LegacyModel bool   `json:"legacy_model"`
	Image       string `json:"image"`
	Phase       string `json:"phase,omitempty"`
	Team        struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
}
