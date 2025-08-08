package repository

import (
	"context"
	"csTrade/internal/domain/offer"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type OfferRepository interface {
	CreateOffer(ctx context.Context, arg offer.OfferDB) error
}

type offerRepository struct {
	db *pgxpool.Pool
}

func NewOfferRepo(db *pgxpool.Pool) OfferRepository {
	return &offerRepository{
		db: db,
	}
}

func (o *offerRepository) CreateOffer(ctx context.Context, arg offer.OfferDB) error {
	query := `
		INSERT INTO offers (
			seller_id, price,
			asset_id, class_id, instance_id, app_id, context_id, amount,
			name, full_name, market_tradable_restriction, icon_url, name_color, action_link,
			tag_type, tag_weapon_internal, tag_weapon_name, tag_quality, tag_rarity, tag_rarity_color, tag_exterior
		)
		VALUES (
			@seller_id, @price,
			@asset_id, @class_id, @instance_id, @app_id, @context_id, @amount,
			@name, @full_name, @market_tradable_restriction, @icon_url, @name_color, @action_link,
			@tag_type, @tag_weapon_internal, @tag_weapon_name, @tag_quality, @tag_rarity, @tag_rarity_color, @tag_exterior
		);
	`

	_, err := o.db.Exec(ctx, query, pgx.NamedArgs{
		"seller_id":                   arg.SellerID,
		"price":                       arg.Price,
		"asset_id":                    arg.AssetID,
		"class_id":                    arg.ClassID,
		"instance_id":                 arg.InstanceID,
		"app_id":                      arg.AppID,
		"context_id":                  arg.ContextID,
		"amount":                      arg.Amount,
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
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateOffer")
		return fmt.Errorf("err create offer db:%w", err)
	}

	return nil
}

func (t *offerRepository) GetOfferByID(ctx context.Context, offerID string) (*offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE id = $1`
	rows, err := t.db.Query(ctx, query, offerID)

	if err != nil {
		return nil, fmt.Errorf("err fetch offers by offer_id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*offer.OfferDB])
}

func (t *offerRepository) GetOfferBySellerID(ctx context.Context, sellerID string) (*offer.OfferDB, error) {
	query := `SELECT * FROM offers WHERE seller_id = $1`
	rows, err := t.db.Query(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("err fetch offers by seller_id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*offer.OfferDB])
}

func (t *offerRepository) UpdateOfferReservedStatus(ctx context.Context, offerID string, reservedTime time.Time) error {

	query := `UPDATE offers SET reserved_until = $1 WHERE id = $2`
	_, err := t.db.Exec(ctx, query, reservedTime, offerID)

	return err
}

func (t *offerRepository) DeleteOfferByID(ctx context.Context, offerID string) error {

	query := `DELETE FROM offers WHERE id = $1`
	_, err := t.db.Exec(ctx, query, offerID)

	return err
}
