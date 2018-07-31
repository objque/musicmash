package services

import (
	"fmt"
	"strings"

	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
)

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

func makeMessage(artist, releaseName, poster string) string {
	return fmt.Sprintf("New release found \n*%s*\n%s [‌‌](%s)", artist, releaseName, poster)
}

func (t *Telegram) Send(args map[string]interface{}) error {
	chatID := args["chatID"].(int64)
	releaseID := args["releaseID"].(uint64)

	release, err := itunes.Lookup(releaseID)
	if err != nil {
		log.Error(errors.Wrapf(err, "can't load information for '%d'", releaseID))
		return err
	}

	message := tgbotapi.NewMessage(chatID, makeMessage(release.ArtistName, release.CollectionName, strings.Replace(release.ArtworkURL100, "100x100", "500x500", 1)))
	message.ParseMode = "markdown"
	message.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Open in iTunes", release.CollectionViewURL)),
		},
	}

	_, err = t.bot.Send(message)
	return err
}
