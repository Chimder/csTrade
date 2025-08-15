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
	CreateUser(ctx context.Context, arg *user.UserCreateReq) error
	// GetUserByID(ctx context.Context, userID string) (*user.UserDB, error)
	GetUserBySteamID(ctx context.Context, steamID string) (*user.UserDB, error)
	GetUserCash(ctx context.Context, userID string) (float64, error)
	GetAllUsers(ctx context.Context) ([]*user.UserDB, error)
	UpdateUserCash(ctx context.Context, cash float64, userID string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (o *userRepository) CreateUser(ctx context.Context, arg *user.UserCreateReq) error {
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

// func (t *userRepository) GetUserByID(ctx context.Context, userID string) (*user.UserDB, error) {
// 	query := `SELECT * FROM users WHERE id = $1`

// 	rows, err := t.db.Query(ctx, query, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("err fetch user by id %w", err)
// 	}

// 	return pgx.CollectOneRow(rows, pgx.RowToStructByName[*user.UserDB])
// }

func (t *userRepository) GetUserBySteamID(ctx context.Context, steamID string) (*user.UserDB, error) {
	log.Info().Str("steamID", steamID).Msg("REPO USER")
	query := `SELECT * FROM users WHERE steam_id = $1`

	rows, err := t.db.Query(ctx, query, steamID)
	if err != nil {
		return nil, fmt.Errorf("err fetch user by steam_id %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.UserDB])
	return &user, err
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
