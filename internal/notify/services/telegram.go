package services

import (
	"fmt"
	"time"

	"github.com/objque/musicmash/internal/clients/itunes/v2"
	"github.com/objque/musicmash/internal/clients/itunes/v2/albums"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	bot      *tgbotapi.BotAPI
	provider *v2.Provider
}

func New(token string, provider *v2.Provider) *Telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Telegram{bot: bot, provider: provider}
}

func makeMessage(release *albums.Album) string {
	releaseDate := ""
	state := "released"
	if release.Attributes.ReleaseDate.Value.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Attributes.ReleaseDate.Value.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.Attributes.Artwork.GetLink(500, 500))
	return fmt.Sprintf("New %s %s \n*%s*\n%s%s %s",
		release.Attributes.GetCollectionType(), state, release.Attributes.ArtistName,
		release.Attributes.Name, releaseDate, poster)
}

func (t *Telegram) Send(args map[string]interface{}) error {
	chatID := args["chatID"].(int64)
	dbRelease := args["release"].(*db.Release)

	release, err := albums.GetAlbumInfo(t.provider, dbRelease.StoreID)
	if err != nil {
		log.Error(errors.Wrapf(err, "can't load information for '%d'", dbRelease.StoreID))
		return err
	}

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
