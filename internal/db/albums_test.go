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

func TestDB_Albums_GetAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: testutil.ReleaseWildwaysTheX}))
	assert.NoError(t, DbMgr.EnsureAlbumExists(&Album{ArtistID: 2, Name: testutil.ReleaseArchitectsHollyHell}))

	// action
	albums, err := DbMgr.GetAlbums(1)

	// assert
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, testutil.ReleaseWildwaysTheX, albums[0].Name)
	assert.Equal(t, uint64(1), albums[0].ID)
}

func TestDB_Albums_GetAlbums_Empty(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: testutil.ReleaseArchitectsHollyHell}))
	assert.NoError(t, DbMgr.EnsureAlbumExists(&Album{ArtistID: 2, Name: testutil.ReleaseWildwaysTheX}))

	// action
	albums, err := DbMgr.GetAlbums(3)

	// assert
	assert.NoError(t, err)
	assert.Empty(t, albums)
}
