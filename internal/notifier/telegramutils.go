package notifier

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func makeText(artistName string, release *db.Release) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.Poster)
	return fmt.Sprintf("New album %s \n*%s*\nby %s%s %s", state, release.Title, artistName, releaseDate, poster)
}

func makeButtons(release *db.Release) *[][]tgbotapi.InlineKeyboardButton {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	storeDetails := config.Config.Stores[release.StoreName]
	buttonLabel := fmt.Sprintf("Listen on %s", storeDetails.Name)
	url := fmt.Sprintf(config.Config.Stores[release.StoreName].ReleaseURL, release.StoreID)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, url)))
	return &buttons
}

func makeMessage(artistName string, release *db.Release) *tgbotapi.MessageConfig {
	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ReplyToMessageID: 0,
			ReplyMarkup:      tgbotapi.InlineKeyboardMarkup{InlineKeyboard: *makeButtons(release)},
		},
		Text:                  makeText(artistName, release),
		ParseMode:             "markdown",
		DisableWebPagePreview: false,
	}
	return &message
}
