package service

import (
	"context"
	"csTrade/internal/domain/offer"
	"csTrade/internal/domain/transaction"
	"csTrade/internal/repository"
	"csTrade/internal/service/bots"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (of *OfferService) ReceiveFromUserOffer(ctx context.Context, offerData *offer.OfferCreateReq) error {
	log.Info().Msg("createOffer")
	bot, errBot := of.botsManager.GetEmptierBot()
	if errBot != nil {
		return errBot
	}

	err := of.repo.WithTxOptions(ctx, pgx.TxOptions{},
		func(r *repository.Repository) error {
			user, err := r.User.GetUserBySteamId(ctx, offerData.SellerID)
			if err != nil {
				return err
			}

			offerData.BotSteamID = bot.SteamID
			offerId, err := r.Offer.CreateOffer(ctx, offerData)
			if offerId == "" {
				return err
			}

			steamTradeId, err := bot.ReceiveFromUser(offerData.AssetID, user.TradeUrl, offerData.SellerID)
			if offerId == "" {
				return err
			}

			err = r.Offer.UpdateOfferAfterReceive(ctx, bot.SteamID, steamTradeId, offerId)
			if err != nil {
				return err
			}

			return nil
		})

	return err
}

func (of *OfferService) GetTradeStatus(ctx context.Context, steamTradeOfferId string) (string, error) {
	offerData, err := of.repo.Offer.GetOfferBySteamOfferID(ctx, steamTradeOfferId)
	if err != nil {
		log.Error().Err(err).Msg("err get offerBotId by steamOfferId")
		return "", fmt.Errorf("err get offerBotId by steamOfferId: %w", err)
	}

	bot := of.botsManager.GetBotByID(offerData.BotSteamID)
	if bot == nil {
		log.Error().Err(err).Msg("err get bot by id")
		return "", fmt.Errorf("err get bot by id")
	}

	err = bot.GetStatus(steamTradeOfferId)
	if err != nil {
		log.Error().Err(err).Msg("err get status by steamOfferId")
		return "", fmt.Errorf("err get status by steamOfferId: %w", err)
	}
	return "ok", err
}

func (of *OfferService) CancelTrade(ctx context.Context, steamTradeOfferId string) error {
	err := of.repo.WithTx(ctx, func(r *repository.Repository) error {
		offerData, err := r.Offer.GetOfferBySteamOfferIDForUpdate(ctx, steamTradeOfferId)
		if err != nil {
			log.Error().Err(err).Msg("err get offerBotId by steamOfferId")
			return fmt.Errorf("err get offerBotId by steamOfferId: %w", err)
		}

		bot := of.botsManager.GetBotByID(offerData.BotSteamID)
		if bot == nil {
			log.Error().Err(err).Msg("err get bot by id")
			return fmt.Errorf("err get bot by id")
		}

		err = bot.DeclineTrade(steamTradeOfferId)
		if err != nil {
			log.Error().Err(err).Msg("err cancel trade")
			return fmt.Errorf("err cancel trade %w", err)
		}

		err = r.Offer.ChangeStatusByID(ctx, offer.OfferCanceled.String(), offerData.ID.String())
		if err != nil {
			log.Error().Err(err).Msg("err change trade statu")
			return fmt.Errorf("err change trade status %w", err)
		}

		return nil
	})

	return err
}

func (of *OfferService) SendToBuyerOffer(ctx context.Context, offer *offer.OfferCreateReq) error {
	log.Info().Msg("sendToBuyerOffer")
	user, err := of.repo.User.GetUserBySteamId(ctx, offer.SellerID)
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

func (of *OfferService) ChangeStatusByID(ctx context.Context, newStatus offer.OfferStatus, offerId string) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("offer status is not valid")
	}

	err := of.repo.Offer.ChangeStatusByID(ctx, newStatus.String(), offerId)
	if err != nil {
		return fmt.Errorf("err change offer status %w", err)
	}

	return err
}
