package transaction

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDB struct {
	ID        uuid.UUID         `db:"id"`
	OfferID   uuid.UUID         `db:"offer_id"`
	SellerID  uint64            `db:"seller_id"`
	BuyerID   uint64            `db:"buyer_id"`
	BotID     uint64            `db:"bot_id"`
	Status    TransactionStatus `db:"status"`
	Price     float64           `db:"price"`
	CreatedAt time.Time         `db:"created_at"`

	Name                      string  `db:"name"`
	FullName                  string  `db:"full_name"`
	MarketTradableRestriction int     `db:"market_tradable_restriction"`
	IconURL                   string  `db:"icon_url"`
	NameColor                 string  `db:"name_color"`
	ActionLink                *string `db:"action_link"`

	TagType           string `db:"tag_type"`
	TagWeaponInternal string `db:"tag_weapon_internal"`
	TagWeaponName     string `db:"tag_weapon_name"`
	TagQuality        string `db:"tag_quality"`
	TagRarity         string `db:"tag_rarity"`
	TagExterior       string `db:"tag_exterior"`
	TagRarityColor    string `db:"tag_rarity_color"`
}

type TransactionStatus string

const (
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
)

func (s TransactionStatus) GetString() string {
	switch s {
	case TransactionFailed:
		return string(TransactionFailed)
	case TransactionCompleted:
		return string(TransactionCompleted)
	}
	return ""
}

func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionCompleted, TransactionFailed:
		return true
	}
	return false
}
