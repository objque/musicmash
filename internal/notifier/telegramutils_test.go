package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTelegramUtils_MakeText(t *testing.T) {
	// arrange
	release := db.Release{
		StoreName: testutils.StoreApple,
		Title:     testutils.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:    testutils.PosterSimple,
	}

	// action
	text := makeText(testutils.ArtistArchitects, &release)

	// assert
	assert.Equal(t, "New album released \n*Holly Hell*\nby Architects [\u200c\u200c](http://pic.jpeg)", text)
}

func TestTelegramUtils_MakeButtons(t *testing.T) {
	// arrange
	release := db.Release{
		StoreName: testutils.StoreApple,
		Title:     testutils.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:    testutils.PosterSimple,
		StoreID:   testutils.StoreIDA,
	}
	config.Config = &config.AppConfig{
		Stores: config.StoresConfig{
			testutils.StoreApple: {
				Name:       testutils.StoreApple,
				ReleaseURL: "https://itunes.apple.com/us/album/%s",
			},
		},
	}

	// action
	buttons := makeButtons(&release)

	// assert
	assert.Len(t, *buttons, 1)
	const wantText = "Listen on " + testutils.StoreApple
	assert.Equal(t, wantText, (*buttons)[0][0].Text)
	const wantURL = "https://itunes.apple.com/us/album/" + testutils.StoreIDA
	assert.NotNil(t, (*buttons)[0][0].URL)
	assert.Equal(t, wantURL, *(*buttons)[0][0].URL)
}

func TestTelegramUtils_MakeText_Announced(t *testing.T) {
	// arrange
	released := time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 48)
	release := db.Release{
		StoreName: testutils.StoreApple,
		Title:     testutils.ReleaseArchitectsHollyHell,
		Released:  released,
		Poster:    testutils.PosterSimple,
	}

	// action
	text := makeText(testutils.ArtistArchitects, &release)

	// assert
	wantMessage := fmt.Sprintf("New album announced \n*Holly Hell*\nby Architects\nRelease date: %s [\u200c\u200c](http://pic.jpeg)", released.Format(time.RFC850))
	assert.Equal(t, wantMessage, text)
}
