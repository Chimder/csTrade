package repository

import (
	"context"
	"csTrade/internal/domain/transaction"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, arg transaction.TransactionDB) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (o *transactionRepository) CreateTransaction(ctx context.Context, arg transaction.TransactionDB) error {
	query := `
		INSERT INTO transactions (
			offer_id, seller_id, buyer_id, status, price, name, full_name, market_tradable_restriction,
			icon_url, name_color, action_link, tag_type, tag_weapon_internal, tag_weapon_name,
			tag_quality, tag_rarity, tag_rarity_color, tag_exterior
		) VALUES (
			@offer_id, @seller_id, @buyer_id, @status, @price, @name, @full_name, @market_tradable_restriction,
			@icon_url, @name_color, @action_link, @tag_type, @tag_weapon_internal, @tag_weapon_name,
			@tag_quality, @tag_rarity, @tag_rarity_color, @tag_exterior
		);
	`

	_, err := o.db.Exec(ctx, query, pgx.NamedArgs{
		"offer_id":                    arg.OfferID,
		"seller_id":                   arg.SellerID,
		"buyer_id":                    arg.BuyerID,
		"status":                      arg.Status,
		"price":                       arg.Price,
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
		log.Error().Err(err).Msg("CreateTransaction")
		return err
	}

	return nil
}
