package notifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func makeText(artistName string, release *db.InternalNotification) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format("Monday, 02-Jan-06"))
	}

	poster := fmt.Sprintf("[\u200c\u200c](%s)", release.Poster)
	return fmt.Sprintf("New %s %s \n*%s*\nby %s%s %s", release.Type, state, release.Title, artistName, releaseDate, poster)
}

func makeButtons(release *db.InternalNotification) *[][]tgbotapi.InlineKeyboardButton {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	storeDetails := config.Config.Stores[release.StoreName]
	releaseAction := "Listen"
	if release.Type == "music-video" {
		releaseAction = "View"
	}
	buttonLabel := fmt.Sprintf("%s on %s", releaseAction, storeDetails.Name)
	url := replacePlaceholders(config.Config.Stores[release.StoreName].ReleaseURL, release)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, url)))
	return &buttons
}

func replacePlaceholders(url string, release *db.InternalNotification) string {
	const (
		replaceCount           = 1
		releaseTypePlaceholder = "{release_type}"
		releaseIDPlaceholder   = "{release_id}"
	)
	url = strings.Replace(url, releaseTypePlaceholder, release.Type, replaceCount)
	return strings.Replace(url, releaseIDPlaceholder, release.StoreID, replaceCount)
}

func makeMessage(artistName string, release *db.InternalNotification) *tgbotapi.MessageConfig {
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
