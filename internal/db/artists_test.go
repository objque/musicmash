package db

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Artist_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExists(testutil.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
}

func TestDB_Artists_GetAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistSkrillex))

	// action
	artists, err := DbMgr.GetAllArtists()

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}

func TestDB_Artists_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistSkrillex))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistSPY))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistWildways))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistRitaOra))
	want := []struct {
		SearchText string
		Artists    []string
	}{
		{SearchText: "il", Artists: []string{testutil.ArtistSkrillex, testutil.ArtistWildways}},
		{SearchText: testutil.ArtistSkrillex, Artists: []string{testutil.ArtistSkrillex}},
		{SearchText: "a", Artists: []string{testutil.ArtistArchitects, testutil.ArtistRitaOra, testutil.ArtistWildways}},
	}

	for _, item := range want {
		// action
		artists, err := DbMgr.SearchArtists(item.SearchText)

		// assert
		assert.NoError(t, err)
		assert.Len(t, artists, len(item.Artists))
		for i, wantName := range item.Artists {
			assert.Equal(t, wantName, artists[i].Name)
		}
	}
}

func TestDB_ArtistStoreInfo_EnsureArtistExistsInStore(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreDeezer, testutil.StoreIDA)

	// assert
	assert.NoError(t, err)
	artists, err := DbMgr.GetArtistsForStore(testutil.StoreDeezer)
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}

func TestDB_ArtistStoreInfo_GetArtistFromStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreApple, testutil.StoreIDB))

	// action
	artists, err := DbMgr.GetArtistFromStore(testutil.ArtistSkrillex, testutil.StoreApple)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
}

func TestDB_ArtistStoreInfo_GetArtistDetails_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	_, err := DbMgr.GetArtistDetails(testutil.ArtistSkrillex)

	// assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDB_ArtistStoreInfo_GetArtistDetails(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistSkrillex))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreDeezer, testutil.StoreIDB))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistSPY, testutil.StoreApple, testutil.StoreIDB))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.ArtistWildways, testutil.StoreYandex, testutil.StoreIDC))
	// recent releases
	now := time.Now().UTC().Truncate(testutil.Day)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   now.Add(-testutil.Day),
	}))
	// another artists (should not be in the result)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week * 2),
	}))
	// announced releases
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Day),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Month),
	}))
	config.Config = &config.AppConfig{
		Stores: map[string]*config.Store{
			testutil.StoreApple:  {ArtistURL: "http://example.com/artist/%s", ReleaseURL: "http://example.com/album/%s"},
			testutil.StoreDeezer: {ArtistURL: "http://example.com/artist/%s", ReleaseURL: "http://example.com/album/%s"},
		},
	}

	// action
	details, err := DbMgr.GetArtistDetails(testutil.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutil.ArtistSkrillex, details.Artist.Name)

	// assert stores
	assert.Len(t, details.Stores, 2)
	// apple music
	assert.Equal(t, testutil.StoreApple, details.Stores[0].StoreName)
	assert.Equal(t, testutil.StoreIDA, details.Stores[0].StoreID)
	assert.Equal(t, testutil.ArtistSkrillex, details.Stores[0].ArtistName)
	// deezer
	assert.Equal(t, testutil.StoreDeezer, details.Stores[1].StoreName)
	assert.Equal(t, testutil.StoreIDB, details.Stores[1].StoreID)
	assert.Equal(t, testutil.ArtistSkrillex, details.Stores[1].ArtistName)

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
