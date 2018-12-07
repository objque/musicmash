package albums

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/clients/deezer"
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
	provider = deezer.NewProvider(server.URL)
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/artist/%d/albums", testutil.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
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
}`, testutil.ReleaseArchitectsHollyHell)))
	})

	// action
	albums, err := GetArtistAlbums(provider, testutil.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, 1045282092, albums[0].ID)
	assert.Equal(t, testutil.ReleaseArchitectsHollyHell, albums[0].Title)
	assert.Equal(t, "2018-08-31", albums[0].Released.Value.Format(testutil.DateYYYYHHMM))
}

func TestClient_GetLatestArtistAlbum(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/artist/13/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "id": 1084871,
      "title": "Keys to the Building",
      "link": "https://www.deezer.com/album/1084871",
      "cover": "https://api.deezer.com/album/1084871/image",
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/fe1e4452bfa11f387129b57a0693479c/1000x1000-000000-80-0-0.jpg",
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
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/1000x1000-000000-80-0-0.jpg",
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
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/92cd4cac9db60794c7f22d2a84e26a29/1000x1000-000000-80-0-0.jpg",
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
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/28d792c5523f083b2981d2827595b7d1/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/28d792c5523f083b2981d2827595b7d1/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/28d792c5523f083b2981d2827595b7d1/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/28d792c5523f083b2981d2827595b7d1/1000x1000-000000-80-0-0.jpg",
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
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/b827a5feb70261cd7341dd2d5bec6fd1/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/b827a5feb70261cd7341dd2d5bec6fd1/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/b827a5feb70261cd7341dd2d5bec6fd1/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/b827a5feb70261cd7341dd2d5bec6fd1/1000x1000-000000-80-0-0.jpg",
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
      "cover_small": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/56x56-000000-80-0-0.jpg",
      "cover_medium": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/250x250-000000-80-0-0.jpg",
      "cover_big": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/500x500-000000-80-0-0.jpg",
      "cover_xl": "https://cdns-images.dzcdn.net/images/cover/8555e5d09e0e491b9e44d0b7fddaa303/1000x1000-000000-80-0-0.jpg",
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
}`, testutil.ReleaseArchitectsHollyHell)))
	})

	// action
	album, err := GetLatestArtistAlbum(provider, 13)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 73607432, album.ID)
	assert.Equal(t, testutil.ReleaseArchitectsHollyHell, album.Title)
	assert.Equal(t, "2018-10-03", album.Released.Value.Format(testutil.DateYYYYHHMM))
}

func TestClient_GetByID(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	url := fmt.Sprintf("/album/%d", testutil.StoreIDQ)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
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
}`, testutil.StoreIDQ, testutil.ReleaseArchitectsHollyHell)))
	})

	// action
	album, err := GetByID(provider, testutil.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutil.StoreIDQ, album.ID)
	assert.Equal(t, testutil.ReleaseArchitectsHollyHell, album.Title)
	assert.NotEmpty(t, album.Poster)
}
