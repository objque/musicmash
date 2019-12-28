package albums

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/musicmash/musicmash/internal/testutils/vars"
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
	provider = deezer.NewProvider(server.URL)
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/artist/%d/albums", vars.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "id": 1045282092,
      "title": "%s",
      "link": "https://www.deezer.com/album/72000342",
      "cover": "https://api.deezer.com/album/72000342/image",
      "cover_small": "https://e-cdns-images.dzcdn.net/images/cover/bf74fc764097630ba58782ae79cfbee6/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://e-cdns-images.dzcdn.net/images/cover/bf74fc764097630ba58782ae79cfbee6/250x250-000000-80-0-0.jpg",
      "cover_big": "https://e-cdns-images.dzcdn.net/images/cover/bf74fc764097630ba58782ae79cfbee6/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://e-cdns-images.dzcdn.net/images/cover/bf74fc764097630ba58782ae79cfbee6/1000x1000-000000-80-0-0.jpg",
      "genre_id": 116,
      "fans": 194659,
      "release_date": "2018-08-31",
      "record_type": "album",
      "tracklist": "https://api.deezer.com/album/72000342/tracks",
      "explicit_lyrics": true,
      "type": "album"
    }
  ],
  "total": 57,
  "next": "https://api.deezer.com/artist/13/albums?limit=1&index=1"
}`, vars.ReleaseArchitectsHollyHell)))
	})

	// action
	albums, err := GetArtistAlbums(provider, vars.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, 1045282092, albums[0].ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, albums[0].Title)
	assert.Equal(t, "2018-08-31", albums[0].Released.Value.Format(vars.DateYYYYHHMM))
}

func TestClient_GetLatestArtistAlbum(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/artist/13/albums", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "id": 1084871,
      "title": "Keys to the Building",
      "link": "https://www.deezer.com/album/1084871",
      "cover": "https://api.deezer.com/album/1084871/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 85,
      "fans": 914,
      "release_date": "2005-02-21",
      "record_type": "album",
      "tracklist": "https://api.deezer.com/album/1084871/tracks",
      "explicit_lyrics": true,
      "type": "album"
    },
    {
      "id": 73607432,
      "title": "%s",
      "link": "https://www.deezer.com/album/73607432",
      "cover": "https://api.deezer.com/album/73607432/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 152,
      "fans": 496,
      "release_date": "2018-10-03",
      "record_type": "single",
      "tracklist": "https://api.deezer.com/album/73607432/tracks",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 72324432,
      "title": "Hereafter (Single)",
      "link": "https://www.deezer.com/album/72324432",
      "cover": "https://api.deezer.com/album/72324432/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 152,
      "fans": 792,
      "release_date": "2018-09-12",
      "record_type": "single",
      "tracklist": "https://api.deezer.com/album/72324432/tracks",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 66896162,
      "title": "Doomsday (Piano Reprise)",
      "link": "https://www.deezer.com/album/66896162",
      "cover": "https://api.deezer.com/album/66896162/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 152,
      "fans": 303,
      "release_date": "2018-07-13",
      "record_type": "single",
      "tracklist": "https://api.deezer.com/album/66896162/tracks",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 46794612,
      "title": "Doomsday (Single)",
      "link": "https://www.deezer.com/album/46794612",
      "cover": "https://api.deezer.com/album/46794612/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 152,
      "fans": 2071,
      "release_date": "2017-09-07",
      "record_type": "single",
      "tracklist": "https://api.deezer.com/album/46794612/tracks",
      "explicit_lyrics": false,
      "type": "album"
    },
    {
      "id": 11745064,
      "title": "Alpha Omega",
      "link": "https://www.deezer.com/album/11745064",
      "cover": "https://api.deezer.com/album/11745064/image",
      "cover_small": "https://dzcdn.net/images/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://dzcdn.net/images/250x250-000000-80-0-0.jpg",
      "cover_big": "https://dzcdn.net/images/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://dzcdn.net/images/1000x1000-000000-80-0-0.jpg",
      "genre_id": 132,
      "fans": 144,
      "release_date": "2012-05-11",
      "record_type": "single",
      "tracklist": "https://api.deezer.com/album/11745064/tracks",
      "explicit_lyrics": false,
      "type": "album"
    }
  ],
  "total": 17
}`, vars.ReleaseArchitectsHollyHell)))
	})

	// action
	album, err := GetLatestArtistAlbum(provider, 13)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 73607432, album.ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, album.Title)
	assert.Equal(t, "2018-10-03", album.Released.Value.Format(vars.DateYYYYHHMM))
}

func TestClient_GetByID(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/album/%d", vars.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "id": %d,
  "title": "%s",
  "upc": "817424019736",
  "link": "https://www.deezer.com/album/76263542",
  "share": "https://www.deezer.com/album/76263542?utm_source=deezer&utm_content=album-76263542&utm_term=0_1540408114&utm_medium=web",
  "cover": "https://api.deezer.com/album/76263542/image",
  "cover_small": "https://e-cdns-images.dzcdn.net/images/cover/3f4983609cbffd22f0d134f9241ed0fb/56x56-000000-80-0-0.jpg",
  "cover_medium": "https://e-cdns-images.dzcdn.net/images/cover/3f4983609cbffd22f0d134f9241ed0fb/250x250-000000-80-0-0.jpg",
  "cover_big": "https://e-cdns-images.dzcdn.net/images/cover/3f4983609cbffd22f0d134f9241ed0fb/500x500-000000-80-0-0.jpg",
  "cover_xl": "https://e-cdns-images.dzcdn.net/images/cover/3f4983609cbffd22f0d134f9241ed0fb/1000x1000-000000-80-0-0.jpg",
  "genre_id": 152
}`, vars.StoreIDQ, vars.ReleaseArchitectsHollyHell)))
	})

	// action
	album, err := GetByID(provider, vars.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, vars.StoreIDQ, album.ID)
	assert.Equal(t, vars.ReleaseArchitectsHollyHell, album.Title)
	assert.NotEmpty(t, album.Poster)
}
