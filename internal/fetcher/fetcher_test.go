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
	mux.HandleFunc("/v1/artists/store/itunes", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(`[{
				"artist_id": %d,
				"name": "itunes",
				"id": "%s"
			}]`, testutil.StoreIDW, testutil.StoreIDA)))
		})
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", testutil.StoreIDA)
	mux.HandleFunc(url, func(w http.ResponseWriter, _ *http.Request) {
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
	fetchFromServices(getServices()).Wait()

	// assert
	releases, err := db.DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
}
