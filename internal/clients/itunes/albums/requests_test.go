package albums

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	mux      *http.ServeMux
	provider *itunes.Provider
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	provider = itunes.NewProvider(server.URL, "82001a6688a941dea1d35f60a7a0f8c3")
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2013-05-28"
      },
      "id": "1045282092"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}
		`))
	})

	// action
	albums, err := GetArtistAlbums(provider, 182821355)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, "1045282092", albums[0].ID)
	assert.Equal(t, "Daybreaker (Deluxe Edition)", albums[0].Attributes.Name)
	assert.Equal(t, "1045635474", albums[1].ID)
	assert.Equal(t, "The Here and Now", albums[1].Attributes.Name)
}

func TestClient_GetLatestArtistAlbum(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2013-05-28"
      },
      "id": "1045282092"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}
		`))
	})

	// action
	album, err := GetLatestArtistAlbum(provider, 182821355)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "1045282092", album.ID)
	assert.Equal(t, "Daybreaker (Deluxe Edition)", album.Attributes.Name)
}
