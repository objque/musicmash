package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/objque/musicmash/internal/clients/itunes/v2/albums"
	"github.com/objque/musicmash/internal/clients/itunes/v2/types"
	"github.com/stretchr/testify/assert"
)

func TestTelegram_MakeMessage_Released(t *testing.T) {
	// arrange
	release := albums.Album{
		Attributes: albums.AlbumAttributes{
			Name:        "Escape - Single",
			ArtistName:  "Gorgon City",
			IsSingle:    true,
			ReleaseDate: types.Time{Value: time.Now().UTC().Truncate(time.Hour).Add(-time.Hour)},
			Artwork: &albums.AlbumArtwork{
				URL: "pic/url",
			},
		},
	}

	// action
	message := makeMessage(&release)

	// assert
	assert.Equal(t, "New Single released \n*Gorgon City*\nEscape - Single [‌‌](pic/url)", message)
}

func TestTelegram_MakeMessage_Future(t *testing.T) {
	// arrange
	release := albums.Album{
		Attributes: albums.AlbumAttributes{
			Name:        "Escape - Single",
			ArtistName:  "Gorgon City",
			IsSingle:    true,
			ReleaseDate: types.Time{time.Now().UTC().Truncate(time.Hour).Add(time.Hour * 24)},
			Artwork: &albums.AlbumArtwork{
				URL: "pic/url",
			},
		},
	}

	// action
	message := makeMessage(&release)

	// assert
	wantMessage := fmt.Sprintf("New Single announced \n*Gorgon City*\nEscape - Single\nRelease date: %s [‌‌](pic/url)", release.Attributes.ReleaseDate.Value.Format(time.RFC850))
	assert.Equal(t, wantMessage, message)
}
