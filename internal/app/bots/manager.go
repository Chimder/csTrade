package bots

import (
	"csTrade/internal/domain/bots"
	"log"
)

type BotManager struct {
	Bots   map[string]*bots.SteamBot
	Events chan interface{}
}

type OfferCreatedEvent struct {
	OfferID  string
	SellerID string
}

type PurchaseRequestedEvent struct {
	TransactionID string
	BuyerID       string
}

func (bm *BotManager) Start() {
	go func() {
		for event := range bm.Events {
			switch e := event.(type) {

			case OfferCreatedEvent:
				log.Print(e)
			case PurchaseRequestedEvent:
			}
		}
	}()
}

func (m *BotManager) InitBots(cfg []bots.SteamBot) error {
	// for _, c := range cfg {
	// bot := bots.NewSteamClient(c.Username, c.AccessToken)
	// if err := bot.Login(); err != nil {
	// 	return err
	// }
	// m.Bots[bot.Username] = bot
	// }
	// return nil
	return nil
}

func (m *BotManager) GetBotByName(name string) *bots.SteamBot {
	for _, b := range m.Bots {
		if b.Username == name {
			return b
		}
	}
	return nil
}
