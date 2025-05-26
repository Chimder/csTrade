package main

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
