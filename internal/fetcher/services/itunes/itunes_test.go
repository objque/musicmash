package itunes

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testAppleMusicClientSuite) TestFetchAndSave() {
	// arrange
	f := Fetcher{Provider: t.provider, FetchWorkers: 5}
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDA)
	t.mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
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
}`, testutil.StoreIDA)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.Association{
		{
			ArtistID:  testutil.StoreIDQ,
			StoreID:   testutil.StoreIDA,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t.T(), testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), 18, releases[0].Released.Day())
	assert.Equal(t.T(), time.July, releases[0].Released.Month())
	assert.Equal(t.T(), 2025, releases[0].Released.Year())
}

func (t *testAppleMusicClientSuite) TestFetchAndSave_AlreadyExists() {
	// arrange
	f := Fetcher{Provider: t.provider, FetchWorkers: 1}
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDB)
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistID:  testutil.StoreIDQ,
		StoreID:   testutil.StoreIDB,
		StoreName: f.GetStoreName(),
	}))
	t.mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
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
}`, testutil.StoreIDB)))
	})

	// action
	wg := sync.WaitGroup{}
	wg.Add(1)
	storeArtists := []*db.Association{
		{
			ArtistID:  testutil.StoreIDQ,
			StoreID:   testutil.StoreIDB,
			StoreName: f.GetStoreName(),
		},
	}
	f.FetchAndSave(&wg, storeArtists)
	wg.Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), releases[0].ArtistID)
	assert.Equal(t.T(), testutil.StoreIDB, releases[0].StoreID)
	assert.Equal(t.T(), 1, releases[0].Released.Day())
	assert.Equal(t.T(), time.January, releases[0].Released.Month())
	assert.Equal(t.T(), 1, releases[0].Released.Year())
}
