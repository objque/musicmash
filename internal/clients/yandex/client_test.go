package yandex

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYandexClient_Auth(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
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
	mux.HandleFunc("/handlers/music-search.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "text": "gorgon city",
    "artists": {
        "items": [{
            "id": 817678,
            "name": "Gorgon City"
        }]
    }
}`))
	})

	// action
	client := Client{httpClient: &http.Client{}, URL: server.URL, Session: &Session{}}
	result, err := client.Search("gordon city")

	// assert
	assert.NoError(t, err)
	assert.Len(t, result.Artists.Items, 1)
	assert.Equal(t, result.Artists.Items[0].ID, 817678)
	assert.Equal(t, result.Artists.Items[0].Name, "Gorgon City")
}

func TestYandexClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/handlers/artist.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
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
	assert.Equal(t, albums[0].ID, 5647716)
	assert.Equal(t, albums[0].Title, "Escape")
	assert.Equal(t, albums[0].ReleaseYear, 2018)
	assert.Equal(t, albums[0].Released, "2018-08-10T00:00:00+03:00")
	// second album
	assert.Equal(t, albums[1].ID, 6564)
	assert.Equal(t, albums[1].Title, "The system")
	assert.Equal(t, albums[1].ReleaseYear, 2017)
	assert.Equal(t, albums[1].Released, "2017-01-10T00:00:00+03:00")
}
