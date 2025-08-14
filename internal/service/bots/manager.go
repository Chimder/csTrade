package bots

import (
	"context"
	"csTrade/internal/domain/bot"
	"csTrade/internal/repository"
	"fmt"
	"sync"
)

type BotManager struct {
	Bots   map[uint64]*bot.SteamBot
	Events chan interface{}
	repo   *repository.Repository
}

func NewBotManager(repo *repository.Repository) *BotManager {
	return &BotManager{
		Bots:   make(map[uint64]*bot.SteamBot),
		Events: make(chan interface{}),
		repo:   repo,
	}
}

type ReceiveFromUserTrade struct {
	// OfferID  string
	// SellerID string
	AssetID  string
	TradeURL string
}

type SendToBuyerEventTrade struct {
	AssetID       string
	BuyerTradeURL string
	// TransactionID string
	BotSteamID uint64
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
			bot := bot.NewSteamClient(&b)
			if err := bot.Login(); err == nil {
				m.Bots[bot.SteamID] = bot
			}
		}(b)
	}
	wg.Wait()
	return nil
}

// func (bm *BotManager) Start(ctx context.Context) {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	bm.InitBots(ctx)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			close(bm.Events)
// 			log.Error().Msg("ctx bots manager cancel")
// 			return

// 		case e := <-bm.Events:
// 			switch ev := e.(type) {
// 			case ReceiveFromUserTrade:
// 				bot := bm.GetEmptierBot()
// 				if bot == nil {
// 					log.Error().Msg("No available bots")
// 					continue
// 				}
// 				if err := bot.ReceiveFromUser(ev.AssetID, ev.TradeURL); err != nil {
// 					log.Error().Err(err).Msg("Err receive from user")
// 				}
// 			case SendToBuyerEventTrade:
// 				bot := bm.GetBotByID(ev.BotSteamID)
// 				if bot == nil {
// 					log.Error().Msg("Bot not found")
// 					continue
// 				}
// 				if err := bot.SendToBuyer(ev.AssetID, ev.BuyerTradeURL); err != nil {
// 					log.Error().Err(err).Msg("Err send to buyer")
// 				}
// 			}
// 		}
// 	}
// }

func (m *BotManager) GetBotByID(steamID uint64) *bot.SteamBot {
	for _, b := range m.Bots {
		if b.SteamID == steamID {
			return b
		}
	}
	return nil
}

func (m *BotManager) GetEmptierBot() *bot.SteamBot {
	var emptier *bot.SteamBot

	for _, v := range m.Bots {
		if emptier == nil || v.SkinCount < emptier.SkinCount {
			emptier = v
		}
	}
	return emptier
}
