package notifier

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"gopkg.in/telegram-bot-api.v4"
)

// TODO (m.kalinin): extract into config
type service struct {
	Name string
	URL  string
}

func (s *service) getURL(id string) string {
	return fmt.Sprintf(s.URL, id)
}

var services = map[string]*service{
	"yandex": &service{Name: "Yandex", URL: "https://music.yandex.ru/artist/%s"},
	"itunes": &service{Name: "Apple Music", URL: "https://itunes.apple.com/us/album/%s"},
}

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
		buttonLabel := fmt.Sprintf("Open in %s", services[store.StoreName].Name)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, services[store.StoreName].getURL(store.StoreID))))
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
