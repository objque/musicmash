package deezer

import (
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_DeezerLinker_Reserve(t *testing.T) {
	task := NewLinker("http://url.mock")

	// action
	task.reserveArtists([]string{testutil.ArtistSkrillex, testutil.ArtistArchitects})

	// assert
	assert.Len(t, task.reservedArtists, 2)
}

func Test_DeezerLinker_Free(t *testing.T) {
	// arrange
	task := NewLinker("http://url.mock")
	artists := []string{testutil.ArtistSkrillex, testutil.ArtistArchitects}
	task.reserveArtists(artists)
	assert.Len(t, task.reservedArtists, 2)

	// action
	task.freeReservedArtists(artists)

	// assert
	assert.Len(t, task.reservedArtists, 0)
}

func Test_DeezerLinker_Search_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	task := NewLinker("http://url.mock")
	artists := []string{testutil.ArtistSkrillex, testutil.ArtistArchitects}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreDeezer, testutil.StoreIDA))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, testutil.StoreDeezer, testutil.StoreIDB))

	// action
	task.SearchArtists(artists)
}

func Test_DeezerLinker_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	task := NewLinker(server.URL)
	mux.HandleFunc("/search/artist", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{
  "data": [
    {
      "id": 525643,
      "name": "Skrillex",
      "link": "https://www.deezer.com/artist/525643",
      "picture": "https://api.deezer.com/artist/525643/image",
      "picture_small": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/56x56-000000-80-0-0.jpg",
      "picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/250x250-000000-80-0-0.jpg",
      "picture_big": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/500x500-000000-80-0-0.jpg",
      "picture_xl": "https://e-cdns-images.dzcdn.net/images/artist/ce7d706ecbba7161ff10e655a00bcc7a/1000x1000-000000-80-0-0.jpg",
      "nb_album": 42,
      "nb_fan": 3445610,
      "radio": true,
      "tracklist": "https://api.deezer.com/artist/525643/top?limit=50",
      "type": "artist"
    }
  ],
  "total": 24,
  "next": "https://api.deezer.com/search/artist?q=skrillex&limit=1&index=1"
}`))
	})

	// action
	task.SearchArtists([]string{testutil.ArtistArchitects})

	// assert
	artists, err := db.DbMgr.GetArtistsForStore(testutil.StoreDeezer)
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}
