package service

import (
	"context"
	"csTrade/internal/domain/offer"
	"csTrade/internal/domain/transaction"
	"csTrade/internal/repository"
	"csTrade/internal/service/bots"
	"fmt"

	"github.com/rs/zerolog/log"
)

type OfferService struct {
	repo        *repository.Repository
	botsManager *bots.BotManager
	// rdb  *redis.Client
}

func NewOfferService(repo *repository.Repository, botsManager *bots.BotManager) *OfferService {
	return &OfferService{repo: repo, botsManager: botsManager}
}

func (of *OfferService) ReceiveFromUserOffer(ctx context.Context, offer *offer.OfferCreateReq) error {
	log.Info().Msg("createOffer")
	user, err := of.repo.User.GetUserBySteamID(ctx, offer.SellerID)
	if err != nil {
		return err
	}

	log.Info().Msg("start get bot")
	bot, err := of.botsManager.GetEmptierBot()
	if err != nil {
		return err
	}

	offer.BotSteamID = bot.SteamID
	offerId, err := of.repo.Offer.CreateOffer(ctx, offer)
	if offerId == "" {
		return err
	}

	steamOfferId, err := bot.ReceiveFromUser(offer.AssetID, user.TradeUrl, offer.SellerID)
	if offerId == "" {
		return err
	}

	err = of.repo.Offer.UpdateOfferAfterReceive(ctx, bot.SteamID, steamOfferId, offerId)
	if err != nil {
		return err
	}

	return nil
}

func (of *OfferService) CancelTrade(ctx context.Context, steamTradeOfferId string) error {
	botId, err := of.repo.Offer.GetOfferBotIdBySteamOfferID(ctx, steamTradeOfferId)
	if err != nil {
		return fmt.Errorf("err get offerBotId by steamOfferId: %w", err)
	}

	bot := of.botsManager.GetBotByID(botId)
	if bot == nil {
		return fmt.Errorf("err get bot by id")
	}

	return bot.DeclineTrade(steamTradeOfferId)
}

func (of *OfferService) SendToBuyerOffer(ctx context.Context, offer *offer.OfferCreateReq) error {
	log.Info().Msg("sendToBuyerOffer")
	user, err := of.repo.User.GetUserBySteamID(ctx, offer.SellerID)
	if err != nil {
		return err
	}

	bot := of.botsManager.GetBotByID(offer.BotSteamID)
	if bot == nil {
		return fmt.Errorf("err get bot by id")
	}

	offer.BotSteamID = bot.SteamID
	err = of.repo.Transaction.CreateTransaction(ctx, transaction.TransactionDB{})
	if err != nil {
		return err
	}

	err = bot.SendToBuyer(offer.AssetID, user.TradeUrl, offer.SellerID)
	if err != nil {
		return err
	}

	return nil
}

func (of *OfferService) GetAllOffers(ctx context.Context) ([]offer.OfferDB, error) {
	return of.repo.Offer.GetAll(ctx)
}

func (of *OfferService) GetByID(ctx context.Context, offerID string) (*offer.OfferDB, error) {
	return of.repo.Offer.GetByID(ctx, offerID)
}

func (of *OfferService) GetUserOffers(ctx context.Context, id string) ([]offer.OfferDB, error) {

	return of.repo.Offer.GetOfferBySellerID(ctx, id)
}

func (of *OfferService) ChangePriceByID(ctx context.Context, id string, newPrice float64) error {
	return of.repo.Offer.ChangePriceByID(ctx, id, newPrice)
}

func (of *OfferService) DeleteByID(ctx context.Context, id string) error {

	return of.repo.Offer.DeleteOfferByID(ctx, id)
}
