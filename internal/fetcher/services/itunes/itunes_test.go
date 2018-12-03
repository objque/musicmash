package itunes

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
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
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Fetching: config.Fetching{
			CountOfSkippedHours: 8,
		},
	}
	provider = itunes.NewProvider(server.URL, "xxx")
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestFetcher_FetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, f.GetStoreName(), "182821355"))
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "158365636"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	f.FetchAndSave(&wg)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "158365636", releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Released.Day())
	assert.Equal(t, time.July, releases[0].Released.Month())
	assert.Equal(t, 2025, releases[0].Released.Year())
}

func TestFetcher_FetchAndSave_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, f.GetStoreName(), "182821355"))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: "158365636", StoreName: f.GetStoreName()}))
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "158365636"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	f.FetchAndSave(&wg)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "158365636", releases[0].StoreID)
	assert.Equal(t, 1, releases[0].Released.Day())
	assert.Equal(t, time.January, releases[0].Released.Month())
	assert.Equal(t, 1, releases[0].Released.Year())
}
