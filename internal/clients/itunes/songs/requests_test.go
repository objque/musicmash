package songs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/testutils/vars"
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
	provider = itunes.NewProvider(server.URL, vars.TokenSimple, time.Minute)
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistSongs(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/v1/catalog/us/artists/%v/songs", vars.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`
{
  "data": [
    {
      "attributes": {
        "name": "%s",
        "releaseDate": "2013-05-28"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "name": "%s",
        "releaseDate": "2012-07-13"
      },
      "id": "%s"
    }
  ]
}`,
			vars.ReleaseArchitectsHollyHell, vars.StoreIDA,
			vars.ReleaseArchitectsHollyHell, vars.StoreIDB)))
	})

	// action
	songs, err := GetArtistSongs(provider, vars.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, songs, 2)
	assert.Equal(t, vars.StoreIDA, songs[0].ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, songs[0].Attributes.Name)
	assert.Equal(t, vars.StoreIDB, songs[1].ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, songs[1].Attributes.Name)
}

func TestClient_GetLatestArtistSong(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/v1/catalog/us/artists/%d/songs", vars.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`
{
  "data": [
    {
      "attributes": {
        "name": "%s",
        "releaseDate": "2025-05-28"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "name": "The Here and Now",
        "releaseDate": "2020-07-13"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`, vars.ReleaseArchitectsHollyHell, vars.StoreIDA, vars.StoreIDB)))
	})

	// action
	songs, err := GetLatestArtistSongs(provider, vars.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, songs, 2)
	assert.Equal(t, vars.StoreIDA, songs[0].ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, songs[0].Attributes.Name)
	assert.Equal(t, vars.StoreIDB, songs[1].ID)
}
