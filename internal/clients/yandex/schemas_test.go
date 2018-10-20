package yandex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYandexClient_ArtistAlbum_GetPosterWithSize(t *testing.T) {
	// arrange
	album := ArtistAlbum{Poster: "avatars.yandex.net/get-music-content/33216/8299493d.a.3796456-1/%%"}
	want := "https://avatars.yandex.net/get-music-content/33216/8299493d.a.3796456-1/400x400"

	// action
	poster := album.GetPosterWithSize(400, 400)

	// assert
	assert.Equal(t, want, poster)
}
