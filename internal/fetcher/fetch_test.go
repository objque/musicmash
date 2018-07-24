package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		StoreURL: server.URL,
	}
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestFetcher_Fetch(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:       "S.P.Y",
		SearchName: "S.P.Y",
	}))
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
          "resultCount": 1,
          "results": [
            {
              "artistName": "S.P.Y",
              "collectionName": "Hospitality: Summer Drum & Bass 2013",
              "releaseDate": "2025-06-03T07:00:00Z"
            }
          ]
        }`))
	})

	// action
	fetch()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
}

func TestFetcher_Internal_IsMustFetch_FirstRun(t *testing.T) {
	// first run means that no records in last_fetches
	setup()
	defer teardown()

	// action
	must := isMustFetch()

	// assert
	assert.True(t, must)
}

func TestFetcher_Internal_IsMustFetch_ReloadApp_AfterFetching(t *testing.T) {
	// fetch was successful and someone restart the app
	setup()
	defer teardown()

	// arrange
	db.DbMgr.SetLastFetch(time.Now().UTC())

	// action
	must := isMustFetch()

	// assert
	assert.False(t, must)
}

func TestFetcher_Internal_IsMustFetch_ReloadApp_AfterOldestFetching(t *testing.T) {
	// fetch was successful some times ago and someone restart the app
	setup()
	defer teardown()

	// arrange
	db.DbMgr.SetLastFetch(time.Now().UTC().Truncate(time.Hour * 48))

	// action
	must := isMustFetch()

	// assert
	assert.True(t, must)
}
