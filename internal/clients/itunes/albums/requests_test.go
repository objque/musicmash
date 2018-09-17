package albums

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v2 "github.com/objque/musicmash/internal/clients/itunes"
	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	mux      *http.ServeMux
	provider *v2.Provider
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	provider = v2.NewProvider(server.URL, "82001a6688a941dea1d35f60a7a0f8c3")
}

func teardown() {
	server.Close()
}

func TestClient_GetArtistAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2013-05-28"
      },
      "id": "1045282092"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}
		`))
	})

	// action
	albums, err := GetArtistAlbums(provider, 182821355)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, "1045282092", albums[0].ID)
	assert.Equal(t, "Daybreaker (Deluxe Edition)", albums[0].Attributes.Name)
	assert.Equal(t, "1045635474", albums[1].ID)
	assert.Equal(t, "The Here and Now", albums[1].Attributes.Name)
}

func TestClient_GetLatestArtistAlbum(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/artists/182821355/albums", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2013-05-28"
      },
      "id": "1045282092"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "isComplete": true,
        "isSingle": false,
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}
		`))
	})

	// action
	album, err := GetLatestArtistAlbum(provider, 182821355)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "1045282092", album.ID)
	assert.Equal(t, "Daybreaker (Deluxe Edition)", album.Attributes.Name)
}

func TestClient_GetInfo(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/albums/1422828208", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "artistName": "K+Lab & Phonik Ops & Jakeda",
        "artwork": {
          "height": 3000,
          "url": "https://is2-ssl.mzstatic.com/image/thumb/Music118/v4/d9/7f/83/d97f8394-acc4-f7db-09aa-6d16391b9040/cover.jpg/{w}x{h}bb.jpeg",
          "width": 3000
        },
		"trackCount": 15,
        "isComplete": false,
        "isSingle": true,
        "name": "Nightmares - Single",
        "releaseDate": "2018-09-14"
      },
      "href": "/v1/catalog/us/albums/1431244851",
      "id": "1422828208",
      "type": "albums"
    }
  ]
}
		`))
	})

	// action
	album, err := GetAlbumInfo(provider, 1422828208)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "1422828208", album.ID)
	assert.Equal(t, "Nightmares - Single", album.Attributes.Name)
	assert.Equal(t, 15, album.Attributes.TrackCount)
	assert.Equal(t, SingleReleaseType, album.Attributes.GetCollectionType())
}

func TestClient_GetAlbumSongs(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/v1/catalog/us/albums/1422828208/tracks", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "data": [
    {
      "attributes": {
        "albumName": "Nightmares",
        "artistName": "Architects",
        "name": "To the Death",
        "previews": [
          {
            "url": "https://cdn.apple/demo-track.m4a"
          }
        ],
        "releaseDate": "2006-05-15"
      },
      "id": "182821357"
    },
    {
      "attributes": {
        "albumName": "Nightmares",
        "artistName": "Architects",
        "name": "You Don't Walk Away from Disemberment",
        "previews": [
          {
            "url": "https://cdn.apple/demo-track.m4a"
          }
        ],
        "releaseDate": "2006-05-15"
      },
      "id": "182821368"
    },
    {
      "attributes": {
        "albumName": "Nightmares",
        "artistName": "Architects",
        "name": "Minesweeper",
        "previews": [
          {
            "url": "https://cdn.apple/demo-tracked2d/mzaf_4114996282890241045.plus.aac.p.m4a"
          }
        ],
        "releaseDate": "2006-05-15"
      },
      "id": "182821616"
    }
  ]
}
		`))
	})

	// action
	songs, err := GetAlbumSongs(provider, 1422828208)

	// assert
	assert.NoError(t, err)
	assert.Len(t, songs, 3)
	assert.Equal(t, "182821357", songs[0].ID)
	assert.Equal(t, "182821368", songs[1].ID)
	assert.Equal(t, "182821616", songs[2].ID)
}
