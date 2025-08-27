package bots

import (
	"context"
	"csTrade/internal/domain/bot"
	"csTrade/internal/repository"
	"fmt"

	"github.com/rs/zerolog/log"
)

type BotManager struct {
	Bots   map[string]*bot.SteamBot
	Events chan interface{}
	repo   *repository.Repository
}

func NewBotManager(repo *repository.Repository) *BotManager {
	return &BotManager{
		Bots:   make(map[string]*bot.SteamBot),
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

func (m *BotManager) InitBots(ctx context.Context) {
	botDB, err := m.repo.Bot.GetBots(ctx)
	if err != nil {
		return
	}

	for _, b := range botDB {
		bot := bot.NewSteamClient(&b)
		log.Info().Str("::", bot.SteamID).Msg("db")

		if err := bot.Login(); err == nil {
			m.Bots[bot.SteamID] = bot
			log.Info().Str("username", bot.Username).Msg("Bot logged in")
		} else {
			log.Error().Err(err).Str("username", bot.Username).Msg("Failed to login bot")
		}
	}

	log.Info().Int("total_bots", len(m.Bots)).Msg("All bots initialized")
}

// func (m *BotManager) InitBots(ctx context.Context) {
// 	botDB, err := m.repo.Bot.GetBots(ctx)
// 	if err != nil {
// 		return
// 	}
// 	log.Info().Interface("BOT", botDB).Msg("bot from db")
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	for _, b := range botDB {
// 		wg.Add(1)
// 		go func(b repository.Bot) {
// 			defer wg.Done()
// 			bot := bot.NewSteamClient(&b)
// 			log.Info().Interface("BOT NEW", bot).Msg("NEW BOT CLIENTR")
// 			if err := bot.Login(); err == nil {
// 				mu.Lock()
// 				m.Bots[bot.SteamID] = bot
// 				mu.Unlock()
// 			} else {
// 				log.Error().Err(err).Str("username", bot.Username).Msg("Failed to login bot")
// 			}
// 		}(b)
// 	}

// 	log.Info().Msg("add bots")
// 	wg.Wait()
// }

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

func (m *BotManager) GetBotByID(steamID string) *bot.SteamBot {
	log.Info().Msg("start bot get")
	for _, b := range m.Bots {
		if b.SteamID == steamID {
			return b
		}
	}
	return nil
}

func (m *BotManager) GetEmptierBot() (*bot.SteamBot, error) {
	if m == nil || len(m.Bots) == 0 {
		return nil, fmt.Errorf("no available bot")
	}

	var emptier *bot.SteamBot
	for _, v := range m.Bots {

		if emptier == nil {
			emptier = v
		} else {

			if v.SkinCount < emptier.SkinCount {
				emptier = v
			}
		}
	}

	if emptier != nil {
		log.Info().
			Str("selected_bot_id", emptier.SteamID).
			Int("skin_count", emptier.SkinCount).
			Msg("Selected emptier bot")
	}

	return emptier, nil
}
