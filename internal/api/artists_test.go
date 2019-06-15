package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/config"
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
		_ = resp.Body.Close()
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
	defer func() { _ = resp.Body.Close() }()

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPI_Artists_GetDetails_ArtistNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	url := fmt.Sprintf("%s/%s/artists/adam tomas moran", server.URL, testutil.UserObjque)
	resp, err := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_Artists_GetDetails_ArtistNotFound_NameWithSpaces(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(testutil.ArtistSkrillex))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreDeezer, testutil.StoreIDB))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistSPY, testutil.StoreApple, testutil.StoreIDB))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistWildways, testutil.StoreYandex, testutil.StoreIDC))
	// recent releases
	now := time.Now().UTC().Truncate(testutil.Day)
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   now.Add(-testutil.Day),
	}))
	// another artists (should not be in the result)
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week),
	}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week * 2),
	}))
	// announced releases
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Day),
	}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Month),
	}))
	config.Config = &config.AppConfig{
		Stores: map[string]*config.Store{
			testutil.StoreApple:  {ArtistURL: "http://example.com/artist/%s", ReleaseURL: "http://example.com/album/%s", Name: testutil.StoreApple},
			testutil.StoreDeezer: {ArtistURL: "http://example.com/artist/%s", ReleaseURL: "http://example.com/album/%s", Name: testutil.StoreDeezer},
		},
	}

	// action
	url := fmt.Sprintf("%s/%s/artists/%s", server.URL, testutil.UserObjque, testutil.ArtistSkrillex)
	resp, err := http.Get(url)
	defer func() { _ = resp.Body.Close() }()
	details := db.ArtistDetails{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&details))

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert stores
	assert.Len(t, details.Stores, 2)
	// apple music
	assert.Equal(t, testutil.StoreApple, details.Stores[0].StoreName)
	assert.Equal(t, testutil.StoreIDA, details.Stores[0].StoreID)
	assert.Empty(t, details.Stores[0].ArtistName)
	// deezer
	assert.Equal(t, testutil.StoreDeezer, details.Stores[1].StoreName)
	assert.Equal(t, testutil.StoreIDB, details.Stores[1].StoreID)
	assert.Empty(t, details.Stores[1].ArtistName)

	// assert recent releases
	assert.Len(t, details.Releases.Recent, 1)
	assert.Len(t, details.Releases.Recent[0].Stores, 1)
	assert.Equal(t, testutil.StoreApple, details.Releases.Recent[0].Stores[0].StoreName)
	assert.Equal(t, testutil.StoreIDA, details.Releases.Recent[0].Stores[0].StoreID)

	// assert announced releases
	assert.Len(t, details.Releases.Announced, 1)
	assert.Len(t, details.Releases.Announced[0].Stores, 2)
	assert.Equal(t, testutil.StoreApple, details.Releases.Announced[0].Stores[0].StoreName)
	assert.Equal(t, testutil.StoreIDC, details.Releases.Announced[0].Stores[0].StoreID)
	assert.Equal(t, testutil.StoreDeezer, details.Releases.Announced[0].Stores[1].StoreName)
	assert.Equal(t, testutil.StoreIDC, details.Releases.Announced[0].Stores[1].StoreID)
}
