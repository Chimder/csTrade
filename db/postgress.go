package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func DBConn(ctx context.Context, url string) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to parse db cfg")
	}


	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 2 * time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed conn to db")
		return nil, err
	}

	log.Ctx(ctx).Info().Str("db_url", url).Msg("Db pool created")
	return pool, nil
}
