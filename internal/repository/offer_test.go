//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"

	"csTrade/internal/domain/offer"
	"csTrade/internal/domain/transaction"
	"csTrade/internal/domain/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:latest",
		tcpostgres.WithDatabase("test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
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

func RandPrice() float64 {
	rawPrice := gofakeit.Float64Range(50, 2000)
	price := math.Round(rawPrice*100) / 100

	return price
}

func GetRandUser(idx int) *user.UserCreateReq {
	userId := fmt.Sprintf("steam_%v%s", idx, gofakeit.DigitN(6))
	return &user.UserCreateReq{
		SteamID:   userId,
		Username:  gofakeit.Username(),
		Email:     gofakeit.Email(),
		Name:      gofakeit.Name(),
		TradeUrl:  "tradeUrl" + gofakeit.Name(),
		AvatarURL: "avatarUrl" + gofakeit.URL(),
	}
}

func GetRandBot(idx int) *Bot {
	botId := fmt.Sprintf("bot_%v%s", idx, gofakeit.DigitN(6))
	return &Bot{
		Username:       gofakeit.Username(),
		Password:       gofakeit.Password(true, true, true, true, true, 10),
		SteamID:        botId,
		SharedSecret:   gofakeit.Regex(`[A-Z0-9]{28}=`),
		SkinCount:      gofakeit.Number(0, 10),
		IdentitySecret: gofakeit.Email(),
		DeviceID:       gofakeit.UUID(),
	}
}

func GetRandOffer(id string) *offer.OfferCreateReq {
	return &offer.OfferCreateReq{
		SellerID:                  id,
		Price:                     RandPrice(),
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

	userRepo := NewUserRepository(db)
	offerRepo := NewOfferRepo(db)
	botsRepo := NewBotsRepo(db)
	transactionRepo := NewTransactionRepo(db)

	var wg sync.WaitGroup

	t.Run("CreateBots", func(t *testing.T) {
		for i := range 20 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				botErr := botsRepo.CreateBots(ctx, GetRandBot(i))
				assert.NoError(t, botErr)
			}(i)
		}
		wg.Wait()
	})

	t.Run("CreateUsers", func(t *testing.T) {
		for i := range 200 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				userErr := userRepo.CreateUser(ctx, GetRandUser(i))
				assert.NoError(t, userErr)
			}(i)
		}
		wg.Wait()
	})

	allBots, userErr := botsRepo.GetBots(ctx)
	require.NoError(t, userErr)
	require.NotEmpty(t, allBots)
	require.GreaterOrEqual(t, len(allBots), 1, "one bot")

	allUsers, userErr := userRepo.GetAllUsers(ctx)
	require.NoError(t, userErr)
	require.NotEmpty(t, allUsers)
	require.GreaterOrEqual(t, len(allUsers), 200, "200 users")

	sellers := allUsers[:100]
	buyers := allUsers[100:]

	t.Run("CreateOffers", func(t *testing.T) {

		for i, u := range sellers {
			wg.Add(1)
			go func(i int, u user.UserDB) {
				defer wg.Done()

				userDb, userErr := userRepo.GetUserBySteamId(ctx, u.SteamID)
				assert.NoError(t, userErr)
				assert.NotEmpty(t, userDb)

				newUserCash := RandPrice()
				userErr = userRepo.UpdateUserCash(ctx, newUserCash, u.SteamID)
				assert.NoError(t, userErr)

				cash, err := userRepo.GetUserCash(ctx, u.SteamID)
				assert.NoError(t, err)
				assert.IsType(t, float64(0), cash)
				assert.Equal(t, newUserCash, cash)

				offerID, err := offerRepo.CreateOffer(ctx, GetRandOffer(u.SteamID))
				assert.NoError(t, err)
				assert.NotEmpty(t, offerID)

				got, err := offerRepo.GetByID(ctx, offerID)
				assert.NoError(t, err)
				assert.NotEmpty(t, got)

				err = offerRepo.UpdateOfferAfterReceive(ctx, allBots[rand.Intn(len(allBots))].SteamID, gofakeit.UUID(), offerID)
				assert.NoError(t, err)

				newPrice := RandPrice()
				err = offerRepo.ChangePriceByID(ctx, offerID, newPrice)
				assert.NoError(t, err)

				offerById, err := offerRepo.GetByID(ctx, offerID)
				assert.NoError(t, err)
				assert.NotEmpty(t, offerById)
				assert.Equal(t, newPrice, offerById.Price)

				assert.NotNil(t, offerById.SteamTradeId)
				offerBySteamOfferID, err := offerRepo.GetOfferBySteamOfferID(ctx, *offerById.SteamTradeId)
				assert.NoError(t, err)
				assert.NotEmpty(t, offerBySteamOfferID)

				offers, err := offerRepo.GetOfferBySellerID(ctx, u.SteamID)
				assert.NoError(t, err)
				assert.NotEmpty(t, offers)
			}(i, u)
		}
		wg.Wait()
	})

	offers, err := offerRepo.GetAll(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, offers)
	t.Run("Create transactions", func(t *testing.T) {
		for i, ofr := range offers {
			wg.Add(1)
			go func(i int, ofr offer.OfferDB) {
				defer wg.Done()
				buyer := buyers[i]

				trErr := transactionRepo.CreateTransaction(ctx, transaction.TransactionDB{
					OfferID:  ofr.ID,
					SellerID: ofr.SellerID,
					BuyerID:  buyer.SteamID,
					BotID:    ofr.BotSteamID,
					Status:   transaction.TransactionCompleted,
					Price:    ofr.Price,
				})
				assert.NoError(t, trErr)

				err := offerRepo.ChangeStatusByID(ctx, offer.OfferSold.String(), ofr.ID.String())
				assert.NoError(t, err)

				updOffer, err := offerRepo.GetByID(ctx, ofr.ID.String())
				assert.NoError(t, err)
				assert.Equal(t, offer.OfferSold, updOffer.Status)
			}(i, ofr)
		}
		wg.Wait()
	})

	transactions, err := transactionRepo.GetAllTransaction()
	require.NoError(t, err)
	require.NotEmpty(t, transactions)
	require.Equal(t, len(transactions), len(buyers))
	require.Equal(t, len(transactions), len(sellers))

	t.Run("TestTransactions", func(t *testing.T) {
		for i, tr := range transactions {
			wg.Add(1)
			go func(i int, tr transaction.TransactionDB) {
				defer wg.Done()

				buyer := buyers[i]
				seller := sellers[i]

				transactionByID, trErr := transactionRepo.GetTransactionByID(ctx, tr.ID.String())
				assert.NoError(t, trErr)
				assert.NotEmpty(t, transactionByID)
				assert.Equal(t, tr.ID, transactionByID.ID)

				transactionByBuyerID, trErr := transactionRepo.GetTransactionByBuyerID(ctx, buyer.SteamID)
				assert.NoError(t, trErr)
				assert.NotEmpty(t, transactionByBuyerID)

				transactionBySellerID, trErr := transactionRepo.GetTransactionBySellerID(ctx, seller.SteamID)
				assert.NoError(t, trErr)
				assert.NotEmpty(t, transactionBySellerID)
			}(i, tr)
		}

		wg.Wait()
	})
	// time.Sleep(5 * time.Minute)
}
