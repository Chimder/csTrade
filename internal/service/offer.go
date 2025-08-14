package service

import (
	"context"
	"csTrade/internal/domain/offer"
	"csTrade/internal/repository"
	"csTrade/internal/service/bots"

	"github.com/rs/zerolog/log"
)

type OfferService struct {
	repo    *repository.Repository
	botsMgr *bots.BotManager
	// rdb  *redis.Client
}

func NewOfferService(repo *repository.Repository) *OfferService {
	return &OfferService{repo: repo}
}

func (of *OfferService) CreateOffer(ctx context.Context, offer *offer.OfferCreateReq) error {
	log.Info().Msg("createOffer")
	user, err := of.repo.User.GetUserBySteamID(ctx, offer.SellerID)
	if err != nil {
		return err
	}

	bot := of.botsMgr.GetEmptierBot()

	offer.BotSteamID = bot.SteamID
	err = of.repo.Offer.CreateOffer(ctx, offer)

	bot.ReceiveFromUser(offer.AssetID, user.TradeUrl)

	return nil
}

// func (s *MangaService) ListMangas(ctx context.Context) ([]byte, error) {

// }
