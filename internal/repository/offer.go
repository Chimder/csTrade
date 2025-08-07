package repository

import (
	"context"
	"csTrade/internal/domain/offer"

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
		return err
	}

	return nil
}

// func (o *offerRepository) GetOffer(ctx context.Context, name string) {
// 	query := `SELECT * FROM "Anime" WHERE name = $1`
// 	rows, err := o.db.Query(ctx, query, name)
// 	if err != nil {
// 		return models.MangaRepo{}, err
// 	}
// }
// 	// return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.MangaRepo])
