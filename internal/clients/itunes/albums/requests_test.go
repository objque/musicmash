package albums

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/testutils"
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
	provider = itunes.NewProvider(server.URL, testutils.TokenSimple, time.Minute)
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/v1/catalog/us/artists/%v/albums", testutils.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`
{
  "data": [
    {
      "attributes": {
        "artistName": "%s",
        "isComplete": true,
        "isSingle": false,
        "name": "%s",
        "releaseDate": "2013-05-28"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "%s",
        "releaseDate": "2012-07-13"
      },
      "id": "%s"
    }
  ]
}`,
			testutils.ArtistArchitects, testutils.ReleaseArchitectsHollyHell, testutils.StoreIDA,
			testutils.ReleaseArchitectsHollyHell, testutils.StoreIDB)))
	})

	// action
	albums, err := GetArtistAlbums(provider, testutils.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, testutils.StoreIDA, albums[0].ID)
	assert.Equal(t, testutils.ReleaseArchitectsHollyHell, albums[0].Attributes.Name)
	assert.Equal(t, testutils.StoreIDB, albums[1].ID)
	assert.Equal(t, testutils.ReleaseArchitectsHollyHell, albums[1].Attributes.Name)
}

func TestClient_GetLatestArtistAlbum(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/v1/catalog/us/artists/%d/albums", testutils.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`
{
  "data": [
    {
      "attributes": {
        "artistName": "%s",
        "isComplete": true,
        "isSingle": false,
        "name": "%s",
        "releaseDate": "2025-05-28"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "The Here and Now",
        "releaseDate": "2020-07-13"
      },
      "id": "%s"
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
}`, testutils.ArtistArchitects, testutils.ReleaseArchitectsHollyHell, testutils.StoreIDA, testutils.StoreIDB)))
	})

	// action
	albums, err := GetLatestArtistAlbums(provider, testutils.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, testutils.StoreIDA, albums[0].ID)
	assert.Equal(t, testutils.ReleaseArchitectsHollyHell, albums[0].Attributes.Name)
	assert.Equal(t, testutils.StoreIDB, albums[1].ID)
}
