package artists

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/testutils"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	mux      *http.ServeMux
	provider *deezer.Provider
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	provider = deezer.NewProvider(server.URL)
}

func teardown() {
	server.Close()
}

func TestClient_SearchArtist(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/search/artist", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
  "data": [
    {
      "id": 525643,
      "name": "Skrillex",
      "link": "https://www.deezer.com/artist/525643",
      "picture": "https://api.deezer.com/artist/525643/image",
      "picture_small": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/56x56-000000-80-0-0.jpg",
      "picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/250x250-000000-80-0-0.jpg",
      "picture_big": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/500x500-000000-80-0-0.jpg",
      "picture_xl": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/1000x1000-000000-80-0-0.jpg",
      "nb_album": 42,
      "nb_fan": 3445610,
      "radio": true,
      "tracklist": "https://api.deezer.com/artist/525643/top?limit=50",
      "type": "artist"
    }
  ],
  "total": 24,
  "next": "https://api.deezer.com/search/artist?q=skrillex&limit=1&index=1"
}`))
	})

	// action
	art, err := SearchArtist(provider, testutils.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 525643, art.ID)
	assert.Equal(t, testutils.ArtistSkrillex, art.Name)
}
