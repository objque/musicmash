package notifier

import (
	"fmt"
	"testing"
	"time"

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
