package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	Username       string
	Password       string
	StreamID       string
	SharedSecret   string
	IdentitySecret string
	DeviceID       string
}

type BotsRepository interface {
	GetBots(ctx context.Context) ([]Bot, error)
}

type botsRepository struct {
	db *pgxpool.Pool
}

func NewBotsRepo(db *pgxpool.Pool) BotsRepository {
	return &botsRepository{
		db: db,
	}
}

func (o *botsRepository) GetBots(ctx context.Context) ([]Bot, error) {
	query := `SELECT * FROM bots`

	rows, err := o.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bots : %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Bot])
}
