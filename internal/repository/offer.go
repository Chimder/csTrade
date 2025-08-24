package repository

import (
	"context"
	"csTrade/internal/domain/offer"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type OfferRepository interface {
	CreateOffer(ctx context.Context, arg *offer.OfferCreateReq) (string, error)
	GetByID(ctx context.Context, offerID string) (*offer.OfferDB, error)
	GetOfferBySellerID(ctx context.Context, sellerID string) ([]offer.OfferDB, error)
	GetAll(ctx context.Context) ([]offer.OfferDB, error)
	AddBotSteamID(ctx context.Context, botSteamId string, offerID string) error
	// UpdateOfferReservedStatus(ctx context.Context, offerID string, reservedTime time.Time) error
	UpdateOfferAfterReceive(ctx context.Context, botSteamId, steamTradeOfferId, offerID string) error
	ChangePriceByID(ctx context.Context, offerID string, newPrice float64) error
	DeleteOfferByID(ctx context.Context, offerID string) error
	GetOfferBotIdBySteamOfferID(ctx context.Context, offerID string) (string, error)
}

type offerRepository struct {
	db *pgxpool.Pool
}

func NewOfferRepo(db *pgxpool.Pool) OfferRepository {
	return &offerRepository{
		db: db,
	}
}

func (o *offerRepository) CreateOffer(ctx context.Context, arg *offer.OfferCreateReq) (string, error) {
	log.Info().Msg("CREate offer DB")
	query := `
		INSERT INTO offers (
			seller_id, price,
			asset_id, class_id, instance_id,
			name, full_name, market_tradable_restriction, icon_url, name_color, action_link,
			tag_type, tag_weapon_internal, tag_weapon_name, tag_quality, tag_rarity, tag_rarity_color, tag_exterior
		)
		VALUES (
			@seller_id, @price,
			@asset_id, @class_id, @instance_id,
			@name, @full_name, @market_tradable_restriction, @icon_url, @name_color, @action_link,
			@tag_type, @tag_weapon_internal, @tag_weapon_name, @tag_quality, @tag_rarity, @tag_rarity_color, @tag_exterior
		) RETURNING id;
	`

	var id string
	err := o.db.QueryRow(ctx, query, pgx.NamedArgs{
		"seller_id":                   arg.SellerID,
		"price":                       arg.Price,
		"asset_id":                    arg.AssetID,
		"class_id":                    arg.ClassID,
		"instance_id":                 arg.InstanceID,
		"name":                        arg.Name,
		"full_name":                   arg.FullName,
		"market_tradable_restriction": arg.MarketTradableRestriction,
		"icon_url":                    arg.IconURL,
		"name_color":                  arg.NameColor,
		"action_link":                 arg.ActionLink,
		"tag_type":                    arg.TagType,
		"tag_weapon_internal":         arg.TagWeaponInternal,
		"tag_weapon_name":             arg.TagWeaponName,
		"tag_quality":                 arg.TagQuality,
		"tag_rarity":                  arg.TagRarity,
		"tag_rarity_color":            arg.TagRarityColor,
		"tag_exterior":                arg.TagExterior,
	}).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("CreateOffer")
		return "", fmt.Errorf("err create offer db:%w", err)
	}

	return id, nil
}

func (t *offerRepository) GetByID(ctx context.Context, offerID string) (*offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE id = $1`
	rows, err := t.db.Query(ctx, query, offerID)

	if err != nil {
		return nil, fmt.Errorf("err fetch offers by offer_id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*offer.OfferDB])
}

func (t *offerRepository) GetAll(ctx context.Context) ([]offer.OfferDB, error) {
	query := `SELECT * FROM offers`
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch all offers %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[offer.OfferDB])
}

func (t *offerRepository) GetOfferBySellerID(ctx context.Context, sellerID string) ([]offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE seller_id = $1`
	rows, err := t.db.Query(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("err fetch offers by seller_id %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[offer.OfferDB])
}

func (t *offerRepository) GetOfferBotIdBySteamOfferID(ctx context.Context, steamOfferID string) (string, error) {
	query := `SELECT bot_steam_id FROM offers WHERE steam_trade_offer_id = $1`
	var botSteamID string
	err := t.db.QueryRow(ctx, query, steamOfferID).Scan(&botSteamID)
	if err != nil {
		return "", fmt.Errorf("err fetch bot_steam_id by steam_trade_offer_id %w", err)
	}
	return botSteamID, nil
}

func (t *offerRepository) UpdateOfferAfterReceive(ctx context.Context, botSteamId, steamTradeOfferId, offerID string) error {
	reservedUntil := time.Now().UTC().Add(15 * time.Minute)

	query := `UPDATE offers SET bot_steam_id = $1, reserved_until = $2, steam_trade_offer_id = $3, updated_at = now() WHERE id = $4`
	_, err := t.db.Exec(ctx, query, botSteamId, reservedUntil, steamTradeOfferId, offerID)

	return err
}

func (t *offerRepository) AddBotSteamID(ctx context.Context, botSteamId string, offerID string) error {
	steamIDUint, err := strconv.ParseUint(botSteamId, 10, 64)
	if err != nil {
		return err
	}
	query := `UPDATE offers SET bot_steam_id = $1 WHERE id = $2`
	_, err = t.db.Exec(ctx, query, steamIDUint, offerID)

	return err
}

func (t *offerRepository) UpdateOfferReservedStatus(ctx context.Context, offerID string, reservedTime time.Time) error {

	query := `UPDATE offers SET reserved_until = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, reservedTime, offerID)

	return err
}

func (t *offerRepository) ChangePriceByID(ctx context.Context, offerID string, newPrice float64) error {
	query := `UPDATE offers SET price = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, newPrice, offerID)

	return err
}
func (t *offerRepository) DeleteOfferByID(ctx context.Context, offerID string) error {

	query := `DELETE FROM offers WHERE id = $1`
	_, err := t.db.Exec(ctx, query, offerID)

	return err
}
