package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/objque/musicmash/internal/clients/itunes"
	"github.com/stretchr/testify/assert"
)

func TestTelegram_MakeMessage_Released(t *testing.T) {
	// arrange
	release := itunes.Release{
		CollectionName: "Escape - Single",
		ArtistName:     "Gorgon City",
		Released:       time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		ArtworkURL100:  "pic/url",
	}

	// action
	message := makeMessage(&release)

	// assert
	assert.Equal(t, "New Single released \n*Gorgon City*\nEscape - Single [‌‌](pic/url)", message)
}

func TestTelegram_MakeMessage_Future(t *testing.T) {
	// arrange
	release := itunes.Release{
		CollectionName: "Escape - Single",
		ArtistName:     "Gorgon City",
		Released:       time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 24),
		ArtworkURL100:  "pic/url",
	}

	// action
	message := makeMessage(&release)

	// assert
	wantMessage := fmt.Sprintf("New Single announced \n*Gorgon City*\nEscape - Single\nRelease date: %s [‌‌](pic/url)", release.Released.Format(time.RFC850))
	assert.Equal(t, wantMessage, message)
}
