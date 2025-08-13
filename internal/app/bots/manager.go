package bots

import (
	"context"
	"csTrade/internal/domain/bots"
	"csTrade/internal/repository"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type BotManager struct {
	Bots   map[string]*bots.SteamBot
	Events chan interface{}
	repo   *repository.Repository
}

func NewBotManager(repo *repository.Repository) *BotManager {
	return &BotManager{
		Bots:   make(map[string]*bots.SteamBot),
		Events: make(chan interface{}),
		repo:   repo,
	}
}

type OfferCreatedEvent struct {
	OfferID  string
	SellerID string
}

type PurchaseRequestedEvent struct {
	TransactionID string
	BuyerID       string
}

func (m *BotManager) InitBots(ctx context.Context) error {
	botDB, err := m.repo.Bot.GetBots(ctx)
	if err != nil {
		return fmt.Errorf("err get bots db: %w", err)
	}
	var wg sync.WaitGroup
	for _, b := range botDB {
		wg.Add(1)
		go func(b repository.Bot) {
			defer wg.Done()
			bot := bots.NewSteamClient(b.Username, b.Password, b.StreamID, b.SharedSecret, b.IdentitySecret, b.DeviceID)
			if err := bot.Login(); err == nil {
				m.Bots[bot.Username] = bot
			}
		}(b)
	}
	wg.Wait()
	return nil
}

func (bm *BotManager) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	bm.InitBots(ctx)
	select {

	case <-ctx.Done():
		close(bm.Events)
		log.Error().Msg("ctx bots manager cancel")
		return
	case e := <-bm.Events:

		switch e.(type) {
		case OfferCreatedEvent:
			// log.Print(e)
		case PurchaseRequestedEvent:
			// log.Print(e)
		}
	}
}

func (m *BotManager) GetBotByName(name string) *bots.SteamBot {
	for _, b := range m.Bots {
		if b.Username == name {
			return b
		}
	}
	return nil
}
