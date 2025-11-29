package repository

import (
	"context"
	"csTrade/internal/domain/user"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type UserStore interface {
	CreateUser(ctx context.Context, arg *user.UserCreateReq) error

	GetUserBySteamId(ctx context.Context, steamID string) (*user.UserDB, error)
	GetUserBySteamIdForUpdate(ctx context.Context, steamID string) (*user.UserDB, error)

	GetUserCash(ctx context.Context, userID string) (float64, error)
	GetUserCashForUpdate(ctx context.Context, userID string) (float64, error)

	GetAllUsers(ctx context.Context) ([]user.UserDB, error)

	UpdateUserCash(ctx context.Context, cash float64, userID string) error
}

type UserRepository struct {
	db Querier
}

func NewUserRepository(db Querier) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (o *UserRepository) CreateUser(ctx context.Context, arg *user.UserCreateReq) error {
	query := `
		INSERT INTO users (steam_id, name, username, email, trade_url, avatar_url)
		VALUES (@steam_id, @name, @username, @email, @trade_url, @avatar_url);
	`

	_, err := o.db.Exec(ctx, query, pgx.NamedArgs{
		"steam_id":   arg.SteamID,
		"name":       arg.Name,
		"username":   arg.Username,
		"email":      arg.Email,
		"trade_url":  arg.TradeUrl,
		"avatar_url": arg.AvatarURL,
	})

	if err != nil {
		log.Error().Err(err).Msg("CreateUserDB")
		return err
	}

	return nil
}

func (t *UserRepository) GetUserBySteamId(ctx context.Context, steamID string) (*user.UserDB, error) {
	// log.Info().Str("steamID", steamID).Msg("REPO USER")
	query := `SELECT * FROM users WHERE steam_id = $1`

	rows, err := t.db.Query(ctx, query, steamID)
	if err != nil {
		return nil, fmt.Errorf("err fetch user by steam_id %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.UserDB])
	return &user, err
}

func (t *UserRepository) GetUserBySteamIdForUpdate(ctx context.Context, steamID string) (*user.UserDB, error) {
	log.Info().Str("steamID", steamID).Msg("REPO USER")
	query := `SELECT * FROM users WHERE steam_id = $1 FOR UPDATE`

	rows, err := t.db.Query(ctx, query, steamID)
	if err != nil {
		return nil, fmt.Errorf("err fetch user by steam_id %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.UserDB])
	return &user, err
}

func (t *UserRepository) GetAllUsers(ctx context.Context) ([]user.UserDB, error) {
	query := `SELECT * FROM users`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("err fetch all users %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.UserDB])
	if err != nil {
		return nil, fmt.Errorf("err collectRows all users %w", err)
	}

	return users, err
}

func (t *UserRepository) GetUserCashForUpdate(ctx context.Context, userID string) (float64, error) {
	query := `SELECT cash FROM users WHERE steam_id = $1 FOR UPDATE`

	var cash float64
	err := t.db.QueryRow(ctx, query, userID).Scan(&cash)
	if err != nil {
		return 0, fmt.Errorf("err fetch user cash by id: %w", err)
	}

	return cash, nil
}

func (t *UserRepository) GetUserCash(ctx context.Context, userID string) (float64, error) {
	query := `SELECT cash FROM users WHERE steam_id = $1`

	var cash float64
	err := t.db.QueryRow(ctx, query, userID).Scan(&cash)
	if err != nil {
		return 0, fmt.Errorf("err fetch user cash by id: %w", err)
	}

	return cash, nil
}

func (t *UserRepository) UpdateUserCash(ctx context.Context, cash float64, userID string) error {
	query := `UPDATE users SET cash = $1 WHERE steam_id = $2`

	_, err := t.db.Exec(ctx, query, cash, userID)
	if err != nil {
		return fmt.Errorf("err fetch update user cash by id %w", err)
	}

	return nil
}
