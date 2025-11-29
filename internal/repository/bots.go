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

type BotsStore interface {
	GetBots(ctx context.Context) ([]Bot, error)
	CreateBots(ctx context.Context, arg *Bot) error
}

type BotsRepository struct {
	db Querier
}

func NewBotsRepo(db Querier) *BotsRepository {
	return &BotsRepository{
		db: db,
	}
}

func (o *BotsRepository) CreateBots(ctx context.Context, arg *Bot) error {
	query := `
		INSERT INTO bots (
			steam_id, username, password, shared_secret, skin_count, identity_secret, device_id
		)
		VALUES (
			@steam_id, @username, @password, @shared_secret, @skin_count, @identity_secret, @device_id
		);
	`
	_, err := o.db.Exec(ctx, query, pgx.NamedArgs{
		"steam_id":        arg.SteamID,
		"username":        arg.Username,
		"password":        arg.Password,
		"shared_secret":   arg.SharedSecret,
		"skin_count":      arg.SkinCount,
		"identity_secret": arg.IdentitySecret,
		"device_id":       arg.DeviceID,
	})
	if err != nil {
		return fmt.Errorf("failed to exec bot : %w", err)
	}

	return nil
}

func (o *BotsRepository) GetBots(ctx context.Context) ([]Bot, error) {
	query := `SELECT * FROM bots`

	rows, err := o.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bots : %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Bot])
}
