package common

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TgClient struct {
	bot      *tgbotapi.BotAPI
	TgnMData TgnMData_Type
}

var TgClientArr []TgClient = make([]TgClient, 0)

func NewTgClient(apiKey string) *TgClient {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		panic(err)
	}

	return &TgClient{
		bot: bot,
	}
}

func (c *TgClient) SendMessage(text string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}
