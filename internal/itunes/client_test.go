package itunes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/objque/musicmash/internal/config"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	config.Config = &config.AppConfig{
		StoreURL: server.URL,
	}
}

func teardown() {
	server.Close()
}

func TestClient_Resolve(t *testing.T) {
	// arrange
	setup()
	defer teardown()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
          "resultCount": 1,
          "results": [
            {
              "artistId": 158365636,
              "artistName": "S.P.Y",
              "artistViewUrl": "https://itunes.apple.com/us/release/s-p-y/158365636?uo=4",
              "artworkUrl100": "https://is5-ssl.mzstatic.com/image/xxx.jpg",
              "artworkUrl30": "https://is5-ssl.mzstatic.com/image/xxx.jpg",
              "artworkUrl60": "https://is5-ssl.mzstatic.com/image/xxx.jpg",
              "collectionArtistId": 4940310,
              "collectionArtistName": "Various Artists",
              "collectionCensoredName": "Hospitality: Summer Drum & Bass 2013",
              "collectionExplicitness": "notExplicit",
              "collectionId": 647213327,
              "collectionName": "Hospitality: Summer Drum & Bass 2013",
              "collectionPrice": 11.99,
              "collectionViewUrl": "https://itunes.apple.com/us/album/one-last-quest/647213327?i=647213403&uo=4",
              "country": "USA",
              "currency": "USD",
              "discCount": 1,
              "discNumber": 1,
              "isStreamable": true,
              "kind": "song",
              "previewUrl": "https://cdn.com/plus.aac.p.m4a",
              "primaryGenreName": "Jungle/Drum'n'bass",
              "releaseDate": "2013-06-03T07:00:00Z",
              "trackCensoredName": "One Last Quest",
              "trackCount": 29,
              "trackExplicitness": "notExplicit",
              "trackId": 647213403,
              "trackName": "One Last Quest",
              "trackNumber": 16,
              "trackPrice": 1.29,
              "trackTimeMillis": 357209,
              "trackViewUrl": "https://itunes.apple.com/us/album/one-last-quest/647213327?i=647213403&uo=4",
              "wrapperType": "track"
            }
          ]
        }`))
	})

	// action
	release, err := GetLatestTrackRelease("S.P.Y")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "S.P.Y", release.ArtistName)
}

func TestClient_GetInfo(t *testing.T) {
	// arrange
	setup()
	defer teardown()

	// action
	release, err := GetArtistInfo("https://itunes.apple.com/us/artist/da-tweekaz/289942206?ign-mpt=uo%3D4")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "https://itunes.apple.com/us/album/j%C3%A4germeister-single/1412554258", release.URL)
	assert.Equal(t, 18, release.Date.Day())
	assert.Equal(t, "July", release.Date.Month().String())
	assert.Equal(t, 2018, release.Date.Year())
}
