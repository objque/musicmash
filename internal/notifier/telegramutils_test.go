package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestTelegramUtils_MakeText(t *testing.T) {
	// arrange
	release := db.Release{
		Title:      "Escape - Single",
		ArtistName: "Gorgon City",
		Released:   time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:     "pic/url",
	}

	// action
	text := makeText(&release)

	// assert
	assert.Equal(t, "New album released \n*Gorgon City*\nEscape - Single [‌‌](pic/url)", text)
}

func TestTelegramUtils_MakeText_Announced(t *testing.T) {
	// arrange
	released := time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 48)
	release := db.Release{
		Title:      "Escape - Single",
		ArtistName: "Gorgon City",
		Released:   released,
		Poster:     "pic/url",
	}

	// action
	text := makeText(&release)

	// assert
	wantMessage := fmt.Sprintf("New album announced \n*Gorgon City*\nEscape - Single\nRelease date: %s [‌‌](pic/url)", released.Format(time.RFC850))
	assert.Equal(t, wantMessage, text)
}
