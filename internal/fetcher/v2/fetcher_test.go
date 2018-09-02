package v2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v2 "github.com/objque/musicmash/internal/clients/itunes"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	ituneshandler "github.com/objque/musicmash/internal/fetcher/handlers/itunes"
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
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Fetching: config.Fetching{
			Workers:                    10,
			CountOfSkippedHoursToFetch: 8,
		},
	}
	provider = v2.NewProvider(server.URL, "xxx")
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestFetcher_FetchAndProcess(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:    "Architects",
		StoreID: 182821355,
	}))
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
        "releaseDate": "2025-07-18"
      },
      "id": "158365636"
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
	f := Fetcher{Provider: provider}
	f.RegisterHandler(&ituneshandler.AppleMusicHandler{})
	err := f.FetchAndProcess()

	// assert
	assert.NoError(t, err)
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, uint64(158365636), releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Date.Day())
	assert.Equal(t, "July", releases[0].Date.Month().String())
	assert.Equal(t, 2025, releases[0].Date.Year())
}

func TestFetcher_FetchAndProcess_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:    "Architects",
		StoreID: 182821355,
	}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		StoreID: 158365636,
	}))
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "158365636"
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
	f := Fetcher{Provider: provider}
	f.RegisterHandler(&ituneshandler.AppleMusicHandler{})
	err := f.FetchAndProcess()

	// assert
	assert.NoError(t, err)
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, uint64(158365636), releases[0].StoreID)
	assert.Equal(t, 1, releases[0].Date.Day())
	assert.Equal(t, "January", releases[0].Date.Month().String())
	assert.Equal(t, 1, releases[0].Date.Year())
}
