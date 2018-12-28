package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Artists_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistSkrillex))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistArchitects))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistSPY))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistWildways))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistRitaOra))
	want := []struct {
		SearchText string
		Artists    []string
	}{
		{SearchText: "il", Artists: []string{testutil.ArtistSkrillex, testutil.ArtistWildways}},
		{SearchText: testutil.ArtistSkrillex, Artists: []string{testutil.ArtistSkrillex}},
		{SearchText: "a", Artists: []string{testutil.ArtistArchitects, testutil.ArtistRitaOra, testutil.ArtistWildways}},
	}

	url := fmt.Sprintf("%s/%s/artists?name=", server.URL, testutil.UserObjque)
	for _, item := range want {
		artists := []db.Artist{}

		// action
		resp, err := http.Get(url + item.SearchText)
		assert.NoError(t, err)
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&artists))

		// assert
		assert.NoError(t, err)
		assert.Len(t, artists, len(item.Artists))
		for i, wantName := range item.Artists {
			assert.Equal(t, wantName, artists[i].Name)
		}
	}
}

func TestAPI_Artists_Search_BadRequest(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))

	// action
	url := fmt.Sprintf("%s/%s/artists?name=", server.URL, testutil.UserObjque)
	resp, err := http.Get(url)

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
