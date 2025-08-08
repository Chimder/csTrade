package repository

import (
	"context"
	"csTrade/internal/domain/user"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg user.UserDB) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (o *userRepository) CreateUser(ctx context.Context, arg user.UserDB) error {
	query := `
		INSERT INTO users (steam_id, username, email, trade_url, avatar_url)
		VALUES (@steam_id, @username, @email, @trade_url, @avatar_url);
	`

	_, err := o.db.Exec(ctx, query, pgx.NamedArgs{
		"steam_id":   arg.SteamID,
		"username":   arg.Username,
		"email":      arg.Email,
		"trade_url":  arg.TradeUrl,
		"avatar_url": arg.AvatarURL,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return err
	}

	return nil
}

func (t *userRepository) GetUserByID(ctx context.Context, userID string) (*user.UserDB, error) {
	query := `SELECT * FROM users WHERE id = $1`

	rows, err := t.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("err fetch user by id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*user.UserDB])
}

func (t *userRepository) GetUserBySteamID(ctx context.Context, steamID string) (*user.UserDB, error) {
	query := `SELECT * FROM users WHERE steam_id = $1`

	rows, err := t.db.Query(ctx, query, steamID)
	if err != nil {
		return nil, fmt.Errorf("err fetch user by steam_id %w", err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*user.UserDB])
}

func (t *userRepository) GetAllUsers(ctx context.Context) ([]*user.UserDB, error) {
	query := `SELECT * FROM users`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch all users %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[*user.UserDB])
}

func (t *userRepository) GetUserCash(ctx context.Context, userID string) (float64, error) {
	query := `SELECT cash FROM users WHERE steam_id = $1`

	var cash float64
	err := t.db.QueryRow(ctx, query, userID).Scan(&cash)
	if err != nil {
		return 0, fmt.Errorf("err fetch user cash by id: %w", err)
	}

	return cash, nil
}

func (t *userRepository) UpdateUserCash(ctx context.Context, cash float64, userID string) error {
	query := `UPDATE users SET cash = $1 WHERE id = $2`

	_, err := t.db.Exec(ctx, query, cash, userID)
	if err != nil {
		return fmt.Errorf("err fetch update user cash by id %w", err)
	}

	return nil
}
