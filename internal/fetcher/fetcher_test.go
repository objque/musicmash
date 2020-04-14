package fetcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testFetcherSuite) TestFetchAndSave() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureAssociationExists(vars.StoreIDW, vars.StoreApple, vars.StoreIDA))
	url := fmt.Sprintf("/v1/catalog/us/artists/%s/albums", vars.StoreIDA)
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
		}`, vars.StoreIDA)))
	})

	// action
	fetchFromServices(getServices()).Wait()

	// assert
	releases, err := db.Mgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(vars.StoreIDW), releases[0].ArtistID)
	assert.Equal(t.T(), vars.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), 18, releases[0].Released.Day())
	assert.Equal(t.T(), time.July, releases[0].Released.Month())
	assert.Equal(t.T(), 2025, releases[0].Released.Year())
}
