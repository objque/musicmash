package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/objque/musicmash/internal/db"
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

func makeMessage(release *itunes.Release) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.ArtworkURL100)
	return fmt.Sprintf("New %s %s \n*%s*\n%s%s %s", release.GetCollectionType(), state, release.ArtistName, release.CollectionName, releaseDate, poster)
}

func (t *Telegram) Send(args map[string]interface{}) error {
	chatID := args["chatID"].(int64)
	dbRelease := args["release"].(*db.Release)

	release, err := itunes.Lookup(dbRelease.StoreID)
	if err != nil {
		log.Error(errors.Wrapf(err, "can't load information for '%d'", dbRelease.StoreID))
		return err
	}

	release.ArtworkURL100 = strings.Replace(release.ArtworkURL100, "100x100", "500x500", 1)
	text := makeMessage(release)
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "markdown"
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, store := range dbRelease.Stores {
		text := fmt.Sprintf("Open in %s", store.GetName())
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(text, store.GetLink())))
	}
	message.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = t.bot.Send(message)
	return err
}
