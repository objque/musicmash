package notifier

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func makeText(release *db.Release) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.Poster)
	return fmt.Sprintf("New album %s \n*%s*\n%s%s %s", state, release.ArtistName, release.Title, releaseDate, poster)
}

func makeButtons(release *db.Release) *[][]tgbotapi.InlineKeyboardButton {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, store := range release.Stores {
		buttonLabel := fmt.Sprintf("Open in %s", config.Config.Stores[store.StoreName].Name)
		url := fmt.Sprintf(config.Config.Stores[store.StoreName].ReleaseURL, store.StoreURL)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, url)))
	}
	return &buttons
}

func MakeMessage(release *db.Release) *tgbotapi.MessageConfig {
	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ReplyToMessageID: 0,
			ReplyMarkup:      tgbotapi.InlineKeyboardMarkup{InlineKeyboard: *makeButtons(release)},
		},
		Text:                  makeText(release),
		ParseMode:             "markdown",
		DisableWebPagePreview: false,
	}
	return &message
}
