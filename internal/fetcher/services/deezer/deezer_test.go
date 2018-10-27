package deezer

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
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
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore("Architects", f.GetStoreName(), "182821355"))
	mux.HandleFunc("/artist/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
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
      "id": 73607432,
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
	assert.Equal(t, "73607432", releases[0].StoreID)
	assert.Equal(t, "Royal Beggars (Single)", releases[0].Title)
	assert.Equal(t, "2020-10-03", releases[0].Released.Format("2006-01-02"))
}

func TestFetcher_FetchAndSave_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore("Architects", f.GetStoreName(), "182821355"))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: "73607432", StoreName: f.GetStoreName()}))
	mux.HandleFunc("/artist/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
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
      "id": 73607432,
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
	assert.Equal(t, "73607432", releases[0].StoreID)
	assert.Equal(t, 1, releases[0].Released.Day())
	assert.Equal(t, "January", releases[0].Released.Month().String())
	assert.Equal(t, 1, releases[0].Released.Year())
}

func TestFetcher_ReFetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	f := Fetcher{Provider: provider, FetchWorkers: 1}
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: "76263542", StoreName: f.GetStoreName()}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: "100054", Poster: "http://pic", StoreName: f.GetStoreName()}))
	mux.HandleFunc("/album/76263542", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "id": 76263542,
  "title": "Pandemonium 2.0",
  "cover_big": "https://e-cdns-images.dzcdn.net/1.jpeg"
}`))
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
	assert.Equal(t, "https://e-cdns-images.dzcdn.net/1.jpeg", releases[0].Poster)
	assert.Equal(t, "http://pic", releases[1].Poster)

}
