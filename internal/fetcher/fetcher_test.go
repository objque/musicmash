package fetcher

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetch_FetchAndSave(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureAssociationExists(testutil.StoreIDW, testutil.StoreApple, testutil.StoreIDA))
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
	assert.Equal(t, int64(testutil.StoreIDW), releases[0].ArtistID)
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, 18, releases[0].Released.Day())
	assert.Equal(t, time.July, releases[0].Released.Month())
	assert.Equal(t, 2025, releases[0].Released.Year())
}
