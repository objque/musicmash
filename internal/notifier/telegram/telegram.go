package telegram

import (
	"net/http"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI

func NewWithClient(token string, client *http.Client) error {
	var err error
	bot, err = tgbotapi.NewBotAPIWithClient(token, client)
	return err
}

func SendMessage(message tgbotapi.Chattable) error {
	_, err := bot.Send(message)
	return err
}
