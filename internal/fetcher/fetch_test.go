package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/notify"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

type MockNotifierService struct{}

func (s *MockNotifierService) Send(args map[string]interface{}) error {
	fmt.Printf("args from notify service: '%v'\n", args)
	return nil
}

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	db.DbMgr = db.NewFakeDatabaseMgr()
	notify.Service = &MockNotifierService{}
	config.Config = &config.AppConfig{
		Store: config.Store{
			URL:    server.URL,
			Region: "us",
		},
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

func TestFetcher_Fetch(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:    "S.P.Y",
		StoreID: 182821355,
	}))
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<section class="l-content-width section section--bordered">
        	<div class="l-row">
          	<div class="l-column small-valign-top small-12 medium-6 large-4">
            	<div class="section__nav">
              	<h2 class="section__headline">Latest Release</h2>
            	</div>
        	<a href="https://itunes.apple.com/us/artist/s-p-y/158365636?uo=4" class="featured-album targeted-link"
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2025" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>
		`))
	})

	// action
	fetch()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, uint64(158365636), releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Date.Day())
	assert.Equal(t, "July", releases[0].Date.Month().String())
	assert.Equal(t, 2025, releases[0].Date.Year())
}

func TestFetcher_Fetch_NoNewReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:    "S.P.Y",
		StoreID: 182821355,
	}))
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<section class="l-content-width section section--bordered">
        	<div class="l-row">
          	<div class="l-column small-valign-top small-12 medium-6 large-4">
            	<div class="section__nav">
              	<h2 class="section__headline">Latest Release</h2>
            	</div>
        	<a href="https://itunes.apple.com/us/artist/s-p-y/158365636?uo=4" class="featured-album targeted-link"
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2017" aria-label="July 18, 2017" class="" >Jul 18, 2017</time>
			</span>
		`))
	})

	// action
	fetch()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 0)
}

func TestFetcher_Fetch_ReleaseAlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{
		Name:    "S.P.Y",
		StoreID: 182821355,
	}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "S.P.Y",
		StoreType:  "itunes",
		StoreID:    158365636,
	}))
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<section class="l-content-width section section--bordered">
        	<div class="l-row">
          	<div class="l-column small-valign-top small-12 medium-6 large-4">
            	<div class="section__nav">
              	<h2 class="section__headline">Latest Release</h2>
            	</div>
        	<a href="https://itunes.apple.com/us/artist/s-p-y/158365636?uo=4" class="featured-album targeted-link"
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2017" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>
		`))
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
	db.DbMgr.SetLastFetch(time.Now().UTC().Add(-time.Hour * 48))

	// action
	must := isMustFetch()

	// assert
	assert.True(t, must)
}
