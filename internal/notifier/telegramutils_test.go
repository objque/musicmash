package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTelegramUtils_MakeText(t *testing.T) {
	// arrange
	release := db.Release{
		Title:      testutil.ReleaseArchitectsHollyHell,
		ArtistName: testutil.ArtistArchitects,
		Released:   time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:     testutil.PosterSimple,
	}

	// action
	text := makeText(&release)

	// assert
	assert.Equal(t, "New album released \n*Architects*\nHolly Hell [‌‌](http://pic.jpeg)", text)
}

func TestTelegramUtils_MakeButtons(t *testing.T) {
	// arrange
	release := db.Release{
		Title:      testutil.ReleaseArchitectsHollyHell,
		ArtistName: testutil.ArtistArchitects,
		Released:   time.Now().UTC().Truncate(time.Hour).Add(-time.Hour),
		Poster:     testutil.PosterSimple,
		Stores: []*db.ReleaseStore{
			{StoreName: testutil.StoreApple, StoreID: testutil.StoreIDA},
		},
	}
	config.Config = &config.AppConfig{
		Stores: map[string]*config.Store{
			testutil.StoreApple: {
				Name:       testutil.StoreApple,
				ReleaseURL: "https://itunes.apple.com/us/album/%s",
			},
		},
	}

	// action
	buttons := makeButtons(&release)

	// assert
	assert.Len(t, *buttons, 1)
	const wantText = "Open in " + testutil.StoreApple
	assert.Equal(t, wantText, (*buttons)[0][0].Text)
	const wantURL = "https://itunes.apple.com/us/album/" + testutil.StoreIDA
	assert.NotNil(t, (*buttons)[0][0].URL)
	assert.Equal(t, wantURL, *(*buttons)[0][0].URL)
}

func TestTelegramUtils_MakeText_Announced(t *testing.T) {
	// arrange
	released := time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 48)
	release := db.Release{
		Title:      testutil.ReleaseArchitectsHollyHell,
		ArtistName: testutil.ArtistArchitects,
		Released:   released,
		Poster:     testutil.PosterSimple,
	}

	// action
	text := makeText(&release)

	// assert
	wantMessage := fmt.Sprintf("New album announced \n*Architects*\nHolly Hell\nRelease date: %s [‌‌](http://pic.jpeg)", released.Format(time.RFC850))
	assert.Equal(t, wantMessage, text)
}
