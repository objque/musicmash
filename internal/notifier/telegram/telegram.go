package telegram

import tgbotapi "gopkg.in/telegram-bot-api.v4"

var bot *tgbotapi.BotAPI

func New(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
}

func SendMessage(message tgbotapi.Chattable) error {
	_, err := bot.Send(message)
	return err
}
