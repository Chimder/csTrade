package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"csTrade/internal/domain/offer"
	"csTrade/internal/domain/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	sqlDB, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	err = goose.Up(sqlDB, "../../sql/migration")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	return pool
}

func GetRandUser(idx int) *user.UserCreateReq {
	userIdx := fmt.Sprintf("steam %v%s", idx, gofakeit.DigitN(6))
	return &user.UserCreateReq{
		SteamID:   userIdx,
		Username:  gofakeit.Username(),
		Email:     gofakeit.Email(),
		Name:      gofakeit.Name(),
		TradeUrl:  "tradeUrl" + gofakeit.Name(),
		AvatarURL: "avatarUrl" + gofakeit.URL(),
	}
}

func GetRandOffer(id string) *offer.OfferCreateReq {
	return &offer.OfferCreateReq{
		SellerID:                  id,
		Price:                     gofakeit.Price(1, 500),
		AssetID:                   gofakeit.UUID(),
		ClassID:                   gofakeit.UUID(),
		InstanceID:                gofakeit.UUID(),
		Name:                      gofakeit.BuzzWord(),
		FullName:                  gofakeit.ProductName(),
		MarketTradableRestriction: gofakeit.Number(1, 14),
		IconURL:                   gofakeit.URL(),
		NameColor:                 gofakeit.Color(),
		ActionLink:                nil,
		TagType:                   "weapon",
		TagWeaponInternal:         gofakeit.Word(),
		TagWeaponName:             gofakeit.CarModel(),
		TagQuality:                gofakeit.Word(),
		TagRarity:                 gofakeit.Word(),
		TagRarityColor:            gofakeit.Color(),
		TagExterior:               gofakeit.RandomString([]string{"FN", "MW", "FT", "WW", "BS"}),
	}
}
func TestOfferRepository(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	gofakeit.Seed(0)

	// transactionRepo := NewTransactionRepo(db)
	userRepo := NewUserRepository(db)
	offerRepo := NewOfferRepo(db)

	var wg sync.WaitGroup
	for i := range 200 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			userErr := userRepo.CreateUser(ctx, GetRandUser(i))
			require.NoError(t, userErr)
		}(i)
	}
	wg.Wait()

	allUsers, userErr := userRepo.GetAllUsers(ctx)
	require.NoError(t, userErr)
	require.NotEmpty(t, allUsers)

	sellers := allUsers[:100]
	buyers := allUsers[100:]

	for i, u := range sellers {
		wg.Add(1)
		go func(i int, u *user.UserDB) {
			defer wg.Done()

			userDb, userErr := userRepo.GetUserBySteamId(ctx, u.SteamID)
			require.NoError(t, userErr)
			require.NotEmpty(t, userDb)

			userErr = userRepo.UpdateUserCash(ctx, gofakeit.Float64(), u.SteamID)
			require.NoError(t, userErr)

			cash, err := userRepo.GetUserCash(ctx, u.SteamID)
			require.NoError(t, userErr)
			require.IsType(t, float64(0), cash)
			////////////////////////////////////////////////////////////

			offerID, err := offerRepo.CreateOffer(ctx, GetRandOffer(u.SteamID))
			require.NoError(t, err)
			require.NotEmpty(t, offerID)

			got, err := offerRepo.GetByID(ctx, offerID)
			require.NoError(t, err)
			require.NotEmpty(t, got)
			offerRepo.
		}(i, u)
	}
	wg.Wait()

	// err = repo.ChangePriceByID(ctx, offerID, 123.45)
	// require.NoError(t, err)

	// updated, err := repo.GetByID(ctx, offerID)
	// require.NoError(t, err)
	// require.Equal(t, 123.45, updated.Price)

	// err = repo.AddBotSteamID(ctx, "987654321", offerID)
	// require.NoError(t, err)

	// afterBot, err := repo.GetByID(ctx, offerID)
	// require.NoError(t, err)
	// require.Equal(t, "987654321", afterBot.BotSteamID)

	// err = repo.ChangeStatusByID(ctx, "reserved", offerID)
	// require.NoError(t, err)

	// afterStatus, err := repo.GetByID(ctx, offerID)
	// require.NoError(t, err)
	// require.Equal(t, "reserved", afterStatus.Status)

	// err = repo.DeleteOfferByID(ctx, offerID)
	// require.NoError(t, err)

	// deleted, err := repo.GetByID(ctx, offerID)
	// require.Error(t, err)
	// require.Nil(t, deleted)
}
