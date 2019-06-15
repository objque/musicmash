package yandex

import (
	"net/http"
	"strings"
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestYandexClient_Auth(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})

	// action
	client := New(server.URL)

	// assert
	assert.Equal(t, "1234276871451297001", client.Session.UID)
}

func TestYandexClient_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/handlers/music-search.jsx", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
    "text": "skrillex",
    "artists": {
        "items": [{
            "id": 817678,
            "name": "Skrillex"
        }]
    }
}`))
	})

	// action
	client := Client{httpClient: &http.Client{}, URL: server.URL, Session: &Session{}}
	result, err := client.Search(strings.ToLower(testutil.ArtistSkrillex))

	// assert
	assert.NoError(t, err)
	assert.Len(t, result.Artists.Items, 1)
	assert.Equal(t, 817678, result.Artists.Items[0].ID)
	assert.Equal(t, testutil.ArtistSkrillex, result.Artists.Items[0].Name)
}

func TestYandexClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/handlers/artist.jsx", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
    "artist": {
        "id": 817678,
        "name": "Gorgon City"
    },
    "albums": [{
        "id": 5647716,
        "title": "Escape",
        "year": 2018,
        "releaseDate": "2018-08-10T00:00:00+03:00"
    },{
        "id": 6564,
        "title": "The system",
        "year": 2017,
        "releaseDate": "2017-01-10T00:00:00+03:00"
    }]
}`))
	})

	// action
	client := Client{httpClient: &http.Client{}, URL: server.URL, Session: &Session{}}
	albums, err := client.GetArtistAlbums(817678)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	// first album
	assert.Equal(t, 5647716, albums[0].ID)
	assert.Equal(t, "Escape", albums[0].Title)
	assert.Equal(t, 2018, albums[0].ReleaseYear)
	assert.Equal(t, "2018-08-10 03:00:00 +0000 UTC", albums[0].Released.Value.String())
	// second album
	assert.Equal(t, 6564, albums[1].ID)
	assert.Equal(t, "The system", albums[1].Title)
	assert.Equal(t, 2017, albums[1].ReleaseYear)
	assert.Equal(t, "2017-01-10 03:00:00 +0000 UTC", albums[1].Released.Value.String())
}
