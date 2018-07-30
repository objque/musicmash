package services

import "gopkg.in/telegram-bot-api.v4"

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func New(token string) *Telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Telegram{bot: bot}
}

func (t *Telegram) Send(args map[string]interface{}) error {
	chatID := args["chatID"].(int64)
	message := args["message"].(string)

	_, err := t.bot.Send(tgbotapi.NewMessage(chatID, message))
	return err
}
