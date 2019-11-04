package itunes

import (
	"fmt"
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
		Fetching: config.FetchingConfig{
			Delay: 8,
		},
	}
	provider = itunes.NewProvider(server.URL, testutil.TokenSimple, time.Minute)
}

func teardown() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
	server.Close()
}

func TestFetcher_FetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 5}
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDA)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
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
      "id": "%s"
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
}`, testutil.StoreIDA)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.ArtistStoreInfo{
		{
			ArtistID:  testutil.StoreIDQ,
			StoreID:   testutil.StoreIDA,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Released.Day())
	assert.Equal(t, time.July, releases[0].Released.Month())
	assert.Equal(t, 2025, releases[0].Released.Year())
}

func TestFetcher_FetchAndSave_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDA)
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistID:  testutil.StoreIDQ,
		StoreID:   testutil.StoreIDB,
		StoreName: f.GetStoreName(),
	}))
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
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
      "id": "%s"
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
}`, testutil.StoreIDB)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.ArtistStoreInfo{
		{
			ArtistID:  testutil.StoreIDQ,
			StoreID:   testutil.StoreIDA,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t, testutil.StoreIDB, releases[0].StoreID)
	assert.Equal(t, 1, releases[0].Released.Day())
	assert.Equal(t, time.January, releases[0].Released.Month())
	assert.Equal(t, 1, releases[0].Released.Year())
}
