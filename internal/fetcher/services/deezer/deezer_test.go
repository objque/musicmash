package deezer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
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
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Fetching: config.Fetching{
			CountOfSkippedHours: 8,
		},
	}
	provider = deezer.NewProvider(server.URL)
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
	url := fmt.Sprintf("/artist/%s/albums", testutil.StoreIDA)
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, f.GetStoreName(), testutil.StoreIDA))
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "id": 1084871,
      "title": "Keys to the Building",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/500x500-000000-80-0-0.jpg",
      "release_date": "2005-02-21",
      "explicit_lyrics": true,
      "type": "album"
    },
    {
      "id": %d,
      "title": "%s",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/500x500-000000-80-0-0.jpg",
      "release_date": "2020-10-03",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 11745064,
      "title": "Alpha Omega",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/500x500-000000-80-0-0.jpg",
      "release_date": "2012-05-11",
      "explicit_lyrics": false,
      "type": "album"
    }
  ],
  "total": 3
}`, testutil.StoreIDQ, testutil.ReleaseArchitectsHollyHell)))
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
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, testutil.ReleaseArchitectsHollyHell, releases[0].Title)
	assert.Equal(t, "2020-10-03", releases[0].Released.Format(testutil.DateYYYYHHMM))
}

func TestFetcher_FetchAndSave_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/artist/%s/albums", testutil.StoreIDA)
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, f.GetStoreName(), testutil.StoreIDA))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: testutil.StoreIDA, StoreName: f.GetStoreName()}))
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "id": 1084871,
      "title": "Keys to the Building",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/500x500-000000-80-0-0.jpg",
      "release_date": "2005-02-21",
      "explicit_lyrics": true,
      "type": "album"
    },
    {
      "id": %s,
      "title": "Royal Beggars (Single)",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/500x500-000000-80-0-0.jpg",
      "release_date": "2020-10-03",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 11745064,
      "title": "Alpha Omega",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/500x500-000000-80-0-0.jpg",
      "release_date": "2012-05-11",
      "explicit_lyrics": false,
      "type": "album"
    }
  ],
  "total": 3
}`, testutil.StoreIDA)))
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
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, 1, releases[0].Released.Day())
	assert.Equal(t, time.January, releases[0].Released.Month())
	assert.Equal(t, 1, releases[0].Released.Year())
}

func TestFetcher_ReFetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/album/%s", testutil.StoreIDA)
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: testutil.StoreIDA, StoreName: f.GetStoreName()}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: testutil.StoreIDB, Poster: testutil.PosterSimple, StoreName: f.GetStoreName()}))
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
  "id": 76263542,
  "title": "Pandemonium 2.0",
  "cover_big": "%s"
}`, testutil.PosterSimple)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	f.ReFetchAndSave(&wg)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, testutil.PosterSimple, releases[0].Poster)
	assert.Equal(t, testutil.PosterSimple, releases[1].Poster)

}
