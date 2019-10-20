package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Albums_IsAlbumExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: testutil.ArtistAlgorithm}))

	// action
	exists := DbMgr.IsAlbumExists(&Album{ArtistID: 1, Name: testutil.ArtistAlgorithm})

	// assert
	assert.True(t, exists)
}

func TestDB_Albums_IsAlbumExists_NotExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	exists := DbMgr.IsAlbumExists(&Album{ArtistID: 1, Name: testutil.ArtistAlgorithm})

	// assert
	assert.False(t, exists)
}
