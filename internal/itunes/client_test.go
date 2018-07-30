package itunes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/objque/musicmash/internal/config"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	config.Config = &config.AppConfig{
		StoreURL: server.URL,
	}
}

func teardown() {
	server.Close()
}

func TestClient_GetInfo(t *testing.T) {
	// arrange
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<section class="l-content-width section section--bordered">
        	<div class="l-row">
          	<div class="l-column small-valign-top small-12 medium-6 large-4">
            	<div class="section__nav">
              	<h2 class="section__headline">Latest Release</h2>
            	</div>
        	<a href="https://itunes.apple.com/us/album/j%C3%A4germeister-single/1412554258" class="featured-album targeted-link"
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2025" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>
		`))
	})

	// action
	release, err := GetArtistInfo(182821355)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "https://itunes.apple.com/us/album/j%C3%A4germeister-single/1412554258", release.URL)
	assert.Equal(t, 18, release.Date.Day())
	assert.Equal(t, "July", release.Date.Month().String())
	assert.Equal(t, 2025, release.Date.Year())
}
