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
	err := DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistSkrillex})

	// assert
	assert.NoError(t, err)
}

func TestDB_Artists_GetAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistSkrillex}))

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
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: testutil.ArtistRitaOra, Followers: 90}))
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

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreDeezer))

	// action
	err := DbMgr.EnsureArtistExistsInStore(testutil.StoreIDQ, testutil.StoreDeezer, testutil.StoreIDA)

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
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreApple))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDB))

	// action
	artists, err := DbMgr.GetArtistFromStore(testutil.StoreIDQ, testutil.StoreApple)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
}

func TestDB_Artists_Validate(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	artists := []int64{1, 5, 10, 420, 69, 333, 999}
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: testutil.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: testutil.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: testutil.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: testutil.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: testutil.ArtistRitaOra, Followers: 90}))

	// action
	artists, err := DbMgr.ValidateArtists(artists)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0])
	assert.Equal(t, int64(5), artists[1])
}

func TestDB_Artists_GetWithFullInfo(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: testutil.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: testutil.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: testutil.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: testutil.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: testutil.ArtistRitaOra, Followers: 90}))

	// action
	artists, err := DbMgr.GetArtistsWithFullInfo([]int64{1, 5})

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)

	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, testutil.ArtistSkrillex, artists[0].Name)
	assert.Equal(t, uint(100), artists[0].Followers)

	assert.Equal(t, int64(5), artists[1].ID)
	assert.Equal(t, testutil.ArtistRitaOra, artists[1].Name)
	assert.Equal(t, uint(90), artists[1].Followers)
}
