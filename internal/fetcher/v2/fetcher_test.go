package v2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers/itunes"
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

func TestFetcher_FetchAndProcess(t *testing.T) {
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
	f := Fetcher{}
	f.RegisterHandler(&itunes.AppleMusicHandler{})
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
		Name:    "S.P.Y",
		StoreID: 182821355,
	}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		StoreType: "itunes",
		StoreID:   158365636,
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
	f := Fetcher{}
	f.RegisterHandler(&itunes.AppleMusicHandler{})
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
