package itunes

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testAppleFetcherSuite) TestFetchAndSave() {
	// arrange
	f := NewService(t.provider, 5, 1)
	albumsURL := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", vars.StoreIDA)
	t.mux.HandleFunc(albumsURL, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`, vars.StoreIDA)))
	})
	songsURL := fmt.Sprintf("/v1/catalog/us/artists/%s/songs", vars.StoreIDA)
	t.mux.HandleFunc(songsURL, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "attributes": {
        "albumName": "The Here and Now",
        "name": "Where Are Ü Now (with Justin Bieber)",
        "releaseDate": "2025-07-18",
		"url": "url:/%s"
      }
    }
  ]
}`, vars.StoreIDB)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.Association{
		{
			ArtistID:  vars.StoreIDQ,
			StoreID:   vars.StoreIDA,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 2)
	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t.T(), vars.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), 18, releases[0].Released.Day())
	assert.Equal(t.T(), time.July, releases[0].Released.Month())
	assert.Equal(t.T(), 2025, releases[0].Released.Year())
	assert.Equal(t.T(), "album", releases[0].Type)

	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[1].ArtistID)
	assert.Equal(t.T(), vars.StoreIDB, releases[1].StoreID)
	assert.Equal(t.T(), 18, releases[1].Released.Day())
	assert.Equal(t.T(), time.July, releases[1].Released.Month())
	assert.Equal(t.T(), 2025, releases[1].Released.Year())
	assert.Equal(t.T(), "song", releases[1].Type)
}

func (t *testAppleFetcherSuite) TestFetchAndSave_AlreadyExists() {
	// arrange
	f := NewService(t.provider, 5, 1)
	albumsURL := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", vars.StoreIDB)
	t.mux.HandleFunc(albumsURL, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "%s"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`, vars.StoreIDA)))
	})
	songsURL := fmt.Sprintf("/v1/catalog/us/artists/%s/songs", vars.StoreIDB)
	t.mux.HandleFunc(songsURL, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`{
  "data": [
    {
      "attributes": {
        "albumName": "The Here and Now",
        "releaseDate": "2025-07-18",
        "url": "url:/%s"
      }
    }
  ]
}`, vars.StoreIDB)))
	})
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistID:  vars.StoreIDQ,
		StoreID:   vars.StoreIDA,
		StoreName: f.GetStoreName(),
	}))
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistID:  vars.StoreIDQ,
		StoreID:   vars.StoreIDB,
		StoreName: f.GetStoreName(),
	}))

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.Association{
		{
			ArtistID:  vars.StoreIDQ,
			StoreID:   vars.StoreIDB,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 2)
	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t.T(), vars.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), 1, releases[0].Released.Day())
	assert.Equal(t.T(), time.January, releases[0].Released.Month())
	assert.Equal(t.T(), 1, releases[0].Released.Year())

	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[1].ArtistID)
	assert.Equal(t.T(), vars.StoreIDB, releases[1].StoreID)
	assert.Equal(t.T(), 1, releases[1].Released.Day())
	assert.Equal(t.T(), time.January, releases[1].Released.Month())
	assert.Equal(t.T(), 1, releases[1].Released.Year())
}
