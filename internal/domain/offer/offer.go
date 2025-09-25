package offer

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type OfferDB struct {
	ID           uuid.UUID `db:"id"`
	SellerID     string    `db:"seller_id"`
	BotSteamID   string    `db:"bot_steam_id"`
	SteamTradeId *string   `db:"steam_trade_id"`
	Price        float64   `db:"price"`

	Status        OfferStatus `db:"status"`
	ReservedUntil *time.Time  `db:"reserved_until"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	AssetID                   string  `db:"asset_id"`
	ClassID                   string  `db:"class_id"`
	InstanceID                string  `db:"instance_id"`
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

type OfferStatus string

const (
	OfferOnSale   OfferStatus = "onsale"
	OfferReserved OfferStatus = "reserved"
	OfferSold     OfferStatus = "sold"
	OfferCanceled OfferStatus = "canceled"
)

var AllOfferStatuses = []OfferStatus{
	OfferOnSale,
	OfferReserved,
	OfferSold,
	OfferCanceled,
}

func (s OfferStatus) String() string {
	return string(s)
}

func (s OfferStatus) IsValid() bool {
	return slices.Contains(AllOfferStatuses, s)
}
