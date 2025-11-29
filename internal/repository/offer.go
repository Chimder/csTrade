package repository

import (
	"context"
	"csTrade/internal/domain/offer"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type OfferStore interface {
	CreateOffer(ctx context.Context, arg *offer.OfferCreateReq) (string, error)
	GetByID(ctx context.Context, offerID string) (*offer.OfferDB, error)
	GetOfferBySellerID(ctx context.Context, sellerID string) ([]offer.OfferDB, error)
	GetAll(ctx context.Context) ([]offer.OfferDB, error)
	AddBotSteamID(ctx context.Context, botSteamId string, offerID string) error
	// UpdateOfferReservedStatus(ctx context.Context, offerID string, reservedTime time.Time) error
	UpdateOfferAfterReceive(ctx context.Context, botSteamId, steamTradeId, offerID string) error
	ChangePriceByID(ctx context.Context, offerID string, newPrice float64) error
	ChangeStatusByID(ctx context.Context, newStatus string, offerId string) error
	GetOfferBySteamOfferID(ctx context.Context, steamTradeID string) (*offer.OfferDB, error)
	GetOfferBySteamOfferIDForUpdate(ctx context.Context, steamTradeID string) (*offer.OfferDB, error)
}

type OfferRepository struct {
	db Querier
}

func NewOfferRepo(db Querier) *OfferRepository {
	return &OfferRepository{
		db: db,
	}
}

func (o *OfferRepository) CreateOffer(ctx context.Context, arg *offer.OfferCreateReq) (string, error) {
	// log.Info().Msg("CREate offer DB")
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

func (t *OfferRepository) GetByID(ctx context.Context, offerID string) (*offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE id = $1`
	rows, err := t.db.Query(ctx, query, offerID)

	if err != nil {
		return nil, fmt.Errorf("err fetch offer by offer_id %w", err)
	}

	offer, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[offer.OfferDB])
	if err != nil {
		return nil, fmt.Errorf("err collect offer by offer_id %w", err)
	}
	return &offer, err
}

func (t *OfferRepository) GetAll(ctx context.Context) ([]offer.OfferDB, error) {
	query := `SELECT * FROM offers`
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch all offers %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[offer.OfferDB])
}

func (t *OfferRepository) GetOfferBySellerID(ctx context.Context, sellerID string) ([]offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE seller_id = $1`
	rows, err := t.db.Query(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("err fetch offers by seller_id %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[offer.OfferDB])
}

func (t *OfferRepository) GetOfferBySteamOfferID(ctx context.Context, steamTradeID string) (*offer.OfferDB, error) {
	// log.Info().Msg("Start get")
	query := `SELECT * FROM offers WHERE steam_trade_id = $1`

	rows, err := t.db.Query(ctx, query, steamTradeID)
	if err != nil {
		return nil, fmt.Errorf("err fetch offer by steam_trade_id %w", err)
	}
	// log.Info().Msg("end get")
	offerData, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[offer.OfferDB])
	if err != nil {
		return nil, fmt.Errorf("err scan offer row: %w", err)
	}

	return &offerData, nil
}

func (t *OfferRepository) GetOfferBySteamOfferIDForUpdate(ctx context.Context, steamTradeID string) (*offer.OfferDB, error) {
	log.Info().Msg("Start get")
	query := `SELECT * FROM offers WHERE steam_trade_id = $1 FOR UPDATE`

	rows, err := t.db.Query(ctx, query, steamTradeID)
	if err != nil {
		return nil, fmt.Errorf("err fetch offer by steam_trade_id %w", err)
	}
	log.Info().Msg("end get")
	offerData, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[offer.OfferDB])
	if err != nil {
		return nil, fmt.Errorf("err scan offer row: %w", err)
	}

	return &offerData, nil
}
func (t *OfferRepository) UpdateOfferAfterReceive(ctx context.Context, botSteamId, steamTradeId, offerID string) error {
	reservedUntil := time.Now().UTC().Add(15 * time.Minute)

	query := `UPDATE offers SET bot_steam_id = $1, reserved_until = $2, steam_trade_id = $3, updated_at = now() WHERE id = $4`
	_, err := t.db.Exec(ctx, query, botSteamId, reservedUntil, steamTradeId, offerID)

	return err
}

func (t *OfferRepository) AddBotSteamID(ctx context.Context, botSteamId string, offerID string) error {
	steamIDUint, err := strconv.ParseUint(botSteamId, 10, 64)
	if err != nil {
		return err
	}
	query := `UPDATE offers SET bot_steam_id = $1 WHERE id = $2`
	_, err = t.db.Exec(ctx, query, steamIDUint, offerID)

	return err
}

func (t *OfferRepository) UpdateOfferReservedStatus(ctx context.Context, offerID string, reservedTime time.Time) error {

	query := `UPDATE offers SET reserved_until = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, reservedTime, offerID)

	return err
}

func (t *OfferRepository) ChangeStatusByID(ctx context.Context, newStatus string, offerId string) error {
	query := `UPDATE offers SET status = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, newStatus, offerId)

	return err
}

func (t *OfferRepository) ChangePriceByID(ctx context.Context, offerID string, newPrice float64) error {
	query := `UPDATE offers SET price = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, newPrice, offerID)

	return err
}
