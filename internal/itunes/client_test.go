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
		Store: config.Store{
			URL:    server.URL,
			Region: "us",
		},
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
        	<a href="https://itunes.apple.com/us/album/j%C3%A4germeister-single/1412554258/" class="featured-album targeted-link"
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2025" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>
		`))
	})

	// action
	release, err := GetArtistInfo(182821355)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, uint64(1412554258), release.ID)
	assert.Equal(t, 18, release.Date.Day())
	assert.Equal(t, "July", release.Date.Month().String())
	assert.Equal(t, 2025, release.Date.Year())
	assert.False(t, release.IsComing)
}

func TestClient_GetInfo_Coming(t *testing.T) {
	// arrange
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/us/artist/182821355", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<div class="section__nav">
                <h2 class="section__headline">Pre-Release</h2>
            </div>
			<section class="l-content-width section section--bordered">
        	<div class="l-row">
          	<div class="l-column small-valign-top small-12 medium-6 large-4">
            	<div class="section__nav">
              		<h2 class="section__headline">Latest Release</h2>
            	</div>
        	<a href="https://itunes.apple.com/us/album/j%C3%A4germeister-single/1412554258/" class="featured-album targeted-link"
			
			<!-- this html-block must be ignored by our fetcher --> 
        	<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2025" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>

			<div class="featured-album__text">
                <span class="featured-album__text__eyebrow targeted-link__target">
                    COMING Aug 24, 2018
                </span>
                <span class="featured-album__text__headline targeted-link__target">
                  Rainier Fog
                </span>
                  <span class="featured-album__text__subcopy targeted-link__target">
                    10 songs
                  </span>
            </div>
		`))
	})

	// action
	release, err := GetArtistInfo(182821355)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, uint64(1412554258), release.ID)
	assert.Equal(t, 24, release.Date.Day())
	assert.Equal(t, "August", release.Date.Month().String())
	assert.Equal(t, 2018, release.Date.Year())
	assert.True(t, release.IsComing)
}

func TestClient_Lookup(t *testing.T) {
	// arrange
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
          "resultCount": 1,
          "results": [
            {
              "wrapperType": "collection",
              "collectionType": "Album",
              "artistId": 646638705,
              "collectionId": 1416749924,
              "amgArtistId": 3118927,
              "artistName": "Party Favor & Baauer",
              "collectionName": "MDR (Remixes) - Single",
              "collectionCensoredName": "MDR (Remixes) - Single",
              "artistViewUrl": "https://itunes.apple.com/us/artist/party-favor/646638705?uo=4",
              "collectionViewUrl": "https://itunes.apple.com/us/album/mdr-remixes-single/1416749924?uo=4",
              "artworkUrl60": "https://is5-ssl.mzstatic.com/image/thumb/Music128/v4/3b/f6/79/3bf6790f-4d61-202b-9f5e-62dd0467c76c/source/60x60bb.jpg",
              "artworkUrl100": "https://is5-ssl.mzstatic.com/image/thumb/Music128/v4/3b/f6/79/3bf6790f-4d61-202b-9f5e-62dd0467c76c/source/100x100bb.jpg",
              "collectionPrice": 2.58,
              "collectionExplicitness": "notExplicit",
              "trackCount": 2,
              "copyright": "℗ 2018 Area 25",
              "country": "USA",
              "currency": "USD",
              "releaseDate": "2018-07-31T07:00:00Z",
              "primaryGenreName": "Dance"
            }
          ]
        }`))
	})

	// action
	release, err := Lookup(1416749924)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 1416749924, release.CollectionID)
	assert.Equal(t, 31, release.Released.Day())
	assert.Equal(t, "July", release.Released.Month().String())
	assert.Equal(t, 2018, release.Released.Year())
}

func TestClient_FindReleaseID(t *testing.T) {
	// arrange
	urls := []string{
		`<a href="https://itunes.apple.com/us/artist/s-p-y/158365636" class="featured-album targeted-link"`,
		`<a href="https://itunes.apple.com/us/artist/s-p-y/158365636?uo=4" class="featured-album targeted-link"`,
		`<a href="https://itunes.apple.com/us/artist/s-p-y/158365636/foo" class="featured-album targeted-link"`,
		`<a href="https://itunes.apple.com/us/artist/s-p-y/158365636/1234" class="featured-album targeted-link"`,
		`<a href="https://itunes.apple.com/us/artist/123/158365636/1234" class="featured-album targeted-link"`,
		`<a href="https://itunes.apple.com/us/artist/spy-13/158365636/1234" class="featured-album targeted-link"`,
	}

	for _, url := range urls {
		// action
		id, err := findReleaseID(url)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, uint64(158365636), *id)
	}
}

func TestClient_FindDate_Coming(t *testing.T) {
	// arrange
	body := `<div class="featured-album__text">
                <span class="featured-album__text__eyebrow targeted-link__target">
                    COMING Aug 24, 2018
                </span>
                <span class="featured-album__text__headline targeted-link__target">
                  Rainier Fog
                </span>
                  <span class="featured-album__text__subcopy targeted-link__target">
                    10 songs
                  </span>
            </div>`

	// action
	date, err := findComingDate(body)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 24, date.Day())
	assert.Equal(t, "August", date.Month().String())
	assert.Equal(t, 2018, date.Year())
}

func TestClient_FindDate_Release(t *testing.T) {
	// arrange
	body := `<span class="featured-album__text__eyebrow targeted-link__target">
				<time data-test-we-datetime datetime="Jul 18, 2025" aria-label="July 18, 2025" class="" >Jul 18, 2025</time>
			</span>`

	// action
	date, err := findReleaseDate(body)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 18, date.Day())
	assert.Equal(t, "July", date.Month().String())
	assert.Equal(t, 2025, date.Year())
}
