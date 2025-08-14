package repository

import (
	"context"
	"csTrade/internal/domain/transaction"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, arg transaction.TransactionDB) error
	GetAllTransaction() ([]transaction.TransactionDB, error)
	GetTransactionByOfferID(ctx context.Context) (*transaction.TransactionDB, error)
	GetTransactionBySellerID(ctx context.Context) ([]transaction.TransactionDB, error)
	GetTransactionByBuyerID(ctx context.Context) ([]transaction.TransactionDB, error)
	GetTransactionByBotID(ctx context.Context) ([]transaction.TransactionDB, error)
	UpdateTransactionStatus(ctx context.Context, status transaction.TransactionStatus, offerID uuid.UUID) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (t *transactionRepository) CreateTransaction(ctx context.Context, arg transaction.TransactionDB) error {
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

	_, err := t.db.Exec(ctx, query, pgx.NamedArgs{
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

func (t *transactionRepository) GetAllTransaction() ([]transaction.TransactionDB, error) {
	ctx := context.Background()
	query := `SELECT * FROM transactions`
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch user stats  %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[transaction.TransactionDB])
}

func (t *transactionRepository) GetTransactionByOfferID(ctx context.Context) (*transaction.TransactionDB, error) {
	query := `SELECT * FROM transactions WHERE offer_id = $1`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch transaction by offer_id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*transaction.TransactionDB])
}

func (t *transactionRepository) GetTransactionBySellerID(ctx context.Context) ([]transaction.TransactionDB, error) {
	query := `SELECT * FROM transactions WHERE seller_id = $1`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch transaction by seller_id %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[transaction.TransactionDB])
}

func (t *transactionRepository) GetTransactionByBuyerID(ctx context.Context) ([]transaction.TransactionDB, error) {
	query := `SELECT * FROM transactions WHERE buyer_id = $1`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch transaction by buyer_id %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[transaction.TransactionDB])
}

func (t *transactionRepository) GetTransactionByBotID(ctx context.Context) ([]transaction.TransactionDB, error) {
	query := `SELECT * FROM transactions WHERE bot_id = $1`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch transaction by bot_id %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[transaction.TransactionDB])
}

func (t *transactionRepository) UpdateTransactionStatus(ctx context.Context, status transaction.TransactionStatus, offerID uuid.UUID) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid transaction status: %s", status)
	}

	query := `UPDATE transactions SET status = $1 WHERE offer_id = $2`
	_, err := t.db.Exec(ctx, query, status, offerID)
	return err
}
