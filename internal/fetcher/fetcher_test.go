package fetcher

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetch_FetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, testutil.StoreYandex, testutil.StoreIDB))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, testutil.StoreApple, testutil.StoreIDA))
	// yandex API mocks
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"yqandexuid": "1234276871451297001"}`))
	})
	mux.HandleFunc("/handlers/artist.jsx", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{
   "artist": {
       "id": %d,
       "name": "Architects"
   },
   "albums": [{
       "id": 5647716,
       "title": "Daybreaker (Deluxe Edition)",
       "year": 2018,
       "releaseDate": "2025-07-18T00:00:00+03:00"
   },{
       "id": 4147713,
       "name": "The Here and Now",
       "year": 2017,
       "releaseDate": "2012-07-13T00:00:00+03:00"
   }]
}`, testutil.StoreIDW)))
	})
	// itunes API mocks
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDA)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{
  "data": [
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "width": 1500,
          "height": 1500,
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "Daybreaker (Deluxe Edition)",
        "releaseDate": "2025-07-18"
      },
      "id": "158365636"
    },
    {
      "attributes": {
        "artistName": "Architects",
        "artwork": {
          "width": 1500,
          "height": 1500,
          "url": "https://is4-ssl.mzstatic.com/image/thumb/Music/0a/90/94/mzi.nyyoiwvs.jpg/{w}x{h}bb.jpeg"
        },
        "name": "The Here and Now",
        "releaseDate": "2012-07-13"
      },
      "id": "1045635474"
    }
  ]
}`))
	})

	// action
	fetchFromServices(getServices()).Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
}
