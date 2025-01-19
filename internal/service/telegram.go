package service

import (
	"almeqapp/config"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot       *tg.BotAPI
	ChannelID string
	EnableLog bool
}

func New(conf *config.Config) (bot *TelegramBot, err error) {
	BotAPI, err := tg.NewBotAPI(conf.TgToken)
	if err != nil {
		return nil, err
	}

	bot = &TelegramBot{
		Bot:       BotAPI,
		ChannelID: conf.ChatId,
		EnableLog: false,
	}

	return
}

func (t *TelegramBot) SendMessage(message string) (err error) {
	msg := tg.NewMessageToChannel(t.ChannelID, message)
	sentMessage, err := t.Bot.Send(msg)
	if err != nil {
		return
	}
	log.Printf("[%v] message sent to channel - [%s]", sentMessage, t.ChannelID)
	return
}
