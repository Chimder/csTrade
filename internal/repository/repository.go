package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	db   Querier
	pool *pgxpool.Pool

	Offer       OfferRepository
	User        UserRepository
	Transaction TransactionRepository
	Bot         BotsRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	r := &Repository{
		db:   pool,
		pool: pool,
	}

	r.Offer = NewOfferRepo(pool)
	r.User = NewUserRepository(pool)
	r.Transaction = NewTransactionRepo(pool)
	r.Bot = NewBotsRepo(pool)

	return r
}

func (r *Repository) newWithTx(tx pgx.Tx) *Repository {
	return &Repository{
		db:          tx,
		pool:        r.pool,
		Offer:       NewOfferRepo(tx),
		User:        NewUserRepository(tx),
		Transaction: NewTransactionRepo(tx),
		Bot:         NewBotsRepo(tx),
	}
}

func (r *Repository) WithTx(ctx context.Context, fn func(*Repository) error) (err error) {
	return r.WithTxOptions(ctx, pgx.TxOptions{}, fn)
}

func (r *Repository) WithTxOptions(ctx context.Context, opts pgx.TxOptions, fn func(*Repository) error) (err error) {
	tx, err := r.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)

			switch v := p.(type) {
			case error:
				log.Error().Err(v).Msg("tx panic recover")
			default:
				log.Error().Any("panic", v).Msg("tx panic recover")
			}

		} else if err != nil {
			_ = tx.Rollback(ctx)
			log.Error().Err(err).Msg("tx panic recover")
		}
	}()

	txRepo := r.newWithTx(tx)
	if r == txRepo {
		return fmt.Errorf("use not tx repo")
	}

	if err = fn(txRepo); err != nil {
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("commit tx err: %w", commitErr)
	}

	return nil
}
