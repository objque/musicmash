package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func TestTelegramUtils_MakeText(t *testing.T) {
	// arrange
	release := db.Release{
		StoreName: vars.StoreApple,
		Title:     vars.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:    vars.PosterSimple,
		Type:      vars.ReleaseTypeAlbum,
	}

	// action
	text := makeText(vars.ArtistArchitects, &release)

	// assert
	assert.Equal(t, "New album released \n*Holly Hell*\nby Architects [\u200c\u200c](http://pic#1.jpeg)", text)
}

func TestTelegramUtils_MakeButtons(t *testing.T) {
	// arrange
	release := db.Release{
		StoreName: vars.StoreApple,
		Title:     vars.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:    vars.PosterSimple,
		StoreID:   vars.StoreIDA,
		Type:      vars.ReleaseTypeAlbum,
	}
	config.Config = &config.AppConfig{
		Stores: config.StoresConfig{
			vars.StoreApple: {
				Name:       vars.StoreApple,
				ReleaseURL: "https://itunes.apple.com/us/album/%s",
			},
		},
	}

	// action
	buttons := makeButtons(&release)

	// assert
	assert.Len(t, *buttons, 1)
	const wantText = "Listen on " + vars.StoreApple
	assert.Equal(t, wantText, (*buttons)[0][0].Text)
	const wantURL = "https://itunes.apple.com/us/album/" + vars.StoreIDA
	assert.NotNil(t, (*buttons)[0][0].URL)
	assert.Equal(t, wantURL, *(*buttons)[0][0].URL)
}

func TestTelegramUtils_MakeText_Announced(t *testing.T) {
	// arrange
	released := time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 48)
	release := db.Release{
		StoreName: vars.StoreApple,
		Title:     vars.ReleaseArchitectsHollyHell,
		Released:  released,
		Poster:    vars.PosterSimple,
		Type:      vars.ReleaseTypeAlbum,
	}

	// action
	text := makeText(vars.ArtistArchitects, &release)

	// assert
	wantMessage := fmt.Sprintf("New album announced \n*Holly Hell*\nby Architects\nRelease date: %s [\u200c\u200c](http://pic#1.jpeg)", released.Format("Monday, 02-Jan-06"))
	assert.Equal(t, wantMessage, text)
}
