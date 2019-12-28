package artists

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

func TestClient_SearchArtist(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/search", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`
{
  "results": {
    "artists": {
      "data": [
        {
          "attributes": {
            "name": "Architects"
          },
          "id": "%s"
        }
      ]
    }
  }
}`, testutils.StoreIDA)))
	})

	// action
	art, err := SearchArtist(provider, testutils.ArtistArchitects)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutils.StoreIDA, art.ID)
	assert.Equal(t, testutils.ArtistArchitects, art.Attributes.Name)
}
