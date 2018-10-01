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