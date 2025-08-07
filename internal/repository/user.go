package repository

////all bots has many item uniq ID
// get item ///

// delete item

//create

//update

import (
	"context"
	"csTrade/internal/domain/user"

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
