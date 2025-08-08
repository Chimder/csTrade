package offer

import (
	"time"

	"github.com/google/uuid"
)

type OfferDB struct {
	ID        uuid.UUID `db:"id"`
	SellerID  uuid.UUID `db:"seller_id"`
	BotID     uuid.UUID `db:"bot_id"`
	Price     float64   `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	ReservedUntil *time.Time `db:"reserved_until"`

	AssetID                   string  `db:"asset_id"`
	ClassID                   string  `db:"class_id"`
	InstanceID                string  `db:"instance_id"`
	AppID                     int     `db:"app_id"`
	ContextID                 string  `db:"context_id"`
	Amount                    int     `db:"amount"`
	Name                      string  `db:"name"`
	FullName                  string  `db:"full_name"`
	MarketTradableRestriction int     `db:"market_tradable_restriction"`
	IconURL                   string  `db:"icon_url"`
	NameColor                 string  `db:"name_color"`
	ActionLink                *string `db:"action_link"`
	TagType                   string  `db:"tag_type"`
	TagWeaponInternal         string  `db:"tag_weapon_internal"`
	TagWeaponName             string  `db:"tag_weapon_name"`
	TagQuality                string  `db:"tag_quality"`
	TagRarity                 string  `db:"tag_rarity"`
	TagRarityColor            string  `db:"tag_rarity_color"`
	TagExterior               string  `db:"tag_exterior"`
}
