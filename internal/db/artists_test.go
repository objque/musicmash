package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Artist_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExists("skrillex")

	// assert
	assert.NoError(t, err)
}

func TestDB_ArtistStoreInfo_EnsureArtistExistsInStore(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExistsInStore("skrillex", "deezer", "xyz")

	// assert
	assert.NoError(t, err)
	artists, err := DbMgr.GetArtistsForStore("deezer")
	assert.Len(t, artists, 1)
}

func TestDB_ArtistStoreInfo_GetArtistFromStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore("skrillex", "itunes", "123"))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore("skrillex", "itunes", "345"))

	// action
	artists, err := DbMgr.GetArtistFromStore("skrillex", "itunes")

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
}
