package artists

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Store: config.Store{
			URL:    server.URL,
			Region: "us",
		},
	}
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestArtists_EnsureExists(t *testing.T) {
	setup()
	defer teardown()
	// arrange
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
          "resultCount": 2,
          "results": [
            {
              "artistId": 3316749924,
              "artistName": "Party Favor & Moderat"
            },
            {
              "artistId": 1416749924,
              "artistName": "Moderat"
            }
          ]
        }`))
	})
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: "Skrillex", StoreID: 001}))
	userArtists := []string{"Skrillex", "modeRAT"}

	// action
	found, notFound := EnsureExists(userArtists)

	// assert
	assert.Len(t, found, 2)
	assert.Len(t, notFound, 0)
	artists, err := db.DbMgr.GetAllArtists()
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, uint64(001), artists[0].StoreID)
	assert.Equal(t, "Skrillex", artists[0].Name)
	assert.Equal(t, uint64(1416749924), artists[1].StoreID)
	assert.Equal(t, "Moderat", artists[1].Name)
}
