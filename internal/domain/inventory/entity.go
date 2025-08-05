package inventory

import "github.com/google/uuid"

type InventoryDB struct {
	ID         uuid.UUID `db:"id"`
	AssetID    string    `db:"asset_id"`
	ClassID    string    `db:"class_id"`
	InstanceID string    `db:"instance_id"`
	AppID      int       `db:"app_id"`
	ContextID  string    `db:"context_id"`
	Amount     int       `db:"amount"`

	MarketHashName            string `db:"market_hash_name"`
	MarketTradableRestriction int    `db:"market_tradable_restriction"`
	IconURL                   string `db:"icon_url"`
	NameColor                 string `db:"name_color"`

	ActionLink *string `db:"action_link"`

	TagType           string `db:"tag_type"`
	TagWeaponInternal string `db:"tag_weapon_internal"`
	TagWeaponName     string `db:"tag_weapon_name"`
	TagQuality        string `db:"tag_quality"`
	TagRarity         string `db:"tag_rarity"`
	TagRarityColor    string `db:"tag_rarity_color"`
	TagExterior       string `db:"tag_exterior"`
}
