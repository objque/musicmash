package telegram

import tgbotapi "gopkg.in/telegram-bot-api.v4"

var bot *tgbotapi.BotAPI

func New(token string) error {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	return err
}

func SendMessage(message tgbotapi.Chattable) error {
	_, err := bot.Send(message)
	return err
}
