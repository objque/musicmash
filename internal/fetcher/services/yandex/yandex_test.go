package yandex

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/musicmash/musicmash/internal/clients/yandex"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
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
		Fetching: config.Fetching{
			Workers:                    10,
			CountOfSkippedHoursToFetch: 8,
		},
	}
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestFetcher_FetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// mock yandex auth
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	mux.HandleFunc("/handlers/artist.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "artist": {
        "id": 817678,
        "name": "Gorgon City"
    },
    "albums": [{
        "id": 5647716,
        "title": "Escape",
        "year": 2018,
        "releaseDate": "2025-07-18T00:00:00+03:00"
    },{
        "id": 6564,
        "title": "The system",
        "year": 2017,
        "releaseDate": "2017-01-10T00:00:00+03:00"
    }]
}`))
	})
	f := Fetcher{API: yandex.New(server.URL)}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore("Gorgon City", f.GetStoreName(), "817678"))

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	f.FetchAndSave(&wg)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "5647716", releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Released.Day())
	assert.Equal(t, "July", releases[0].Released.Month().String())
	assert.Equal(t, 2025, releases[0].Released.Year())
}

func TestFetcher_FetchAndSave_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// mock yandex auth
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	mux.HandleFunc("/handlers/artist.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "artist": {
        "id": 817678,
        "name": "Gorgon City"
    },
    "albums": [{
        "id": 5647716,
        "title": "Escape",
        "year": 2018,
        "releaseDate": "2025-07-18T00:00:00+03:00"
    },{
        "id": 6564,
        "title": "The system",
        "year": 2017,
        "releaseDate": "2017-01-10T00:00:00+03:00"
    }]
}`))
	})
	f := Fetcher{API: yandex.New(server.URL)}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore("Gorgon City", f.GetStoreName(), "817678"))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{StoreID: "5647716", StoreName: f.GetStoreName()}))

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	f.FetchAndSave(&wg)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "5647716", releases[0].StoreID)
	// NOTE (m.kalinin): mock was created with zero date
	assert.Equal(t, 1, releases[0].Released.Day())
	assert.Equal(t, "January", releases[0].Released.Month().String())
	assert.Equal(t, 1, releases[0].Released.Year())
}
