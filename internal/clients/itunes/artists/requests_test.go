package artists

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v2 "github.com/objque/musicmash/internal/clients/itunes"
	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	mux      *http.ServeMux
	provider *v2.Provider
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	provider = v2.NewProvider(server.URL, "82001a6688a941dea1d35f60a7a0f8c3")
}

func teardown() {
	server.Close()
}

func TestClient_SearchArtist(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "results": {
    "artists": {
      "data": [
        {
          "attributes": {
            "name": "Architects"
          },
          "id": "182821355"
        }
      ]
    }
  }
}
		`))
	})

	// action
	art, err := SearchArtist(provider, "Architects")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "182821355", art.ID)
	assert.Equal(t, "Architects", art.Attributes.Name)
}
