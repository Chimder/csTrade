package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Bot struct {
	Username       string `db:"username"`
	Password       string `db:"password"`
	SteamID        string `db:"steam_id"`
	SharedSecret   string `db:"shared_secret"`
	SkinCount      int    `db:"skin_count"`
	IdentitySecret string `db:"identity_secret"`
	DeviceID       string `db:"device_id"`
}

type BotsRepository interface {
	GetBots(ctx context.Context) ([]Bot, error)
}

type botsRepository struct {
	db Querier
}

func NewBotsRepo(db Querier) BotsRepository {
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
