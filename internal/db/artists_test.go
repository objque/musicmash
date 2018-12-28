package db

import (
	"testing"

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
