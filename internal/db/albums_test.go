package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestAlbums_IsAlbumExists() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: vars.ArtistAlgorithm}))

	// action
	exists := DbMgr.IsAlbumExists(&Album{ArtistID: 1, Name: vars.ArtistAlgorithm})

	// assert
	assert.True(t.T(), exists)
}
func (t *testDBSuite) TestAlbums_IsAlbumExists_NotExists() {
	// action
	exists := DbMgr.IsAlbumExists(&Album{ArtistID: 1, Name: vars.ArtistAlgorithm})

	// assert
	assert.False(t.T(), exists)
}

func (t *testDBSuite) TestAlbums_GetAlbums() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: vars.ReleaseWildwaysTheX}))
	assert.NoError(t.T(), DbMgr.EnsureAlbumExists(&Album{ArtistID: 2, Name: vars.ReleaseArchitectsHollyHell}))

	// action
	albums, err := DbMgr.GetAlbums(1)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), albums, 1)
	assert.Equal(t.T(), vars.ReleaseWildwaysTheX, albums[0].Name)
	assert.Equal(t.T(), uint64(1), albums[0].ID)
}

func (t *testDBSuite) TestAlbums_GetAlbums_Empty() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureAlbumExists(&Album{ArtistID: 1, Name: vars.ReleaseArchitectsHollyHell}))
	assert.NoError(t.T(), DbMgr.EnsureAlbumExists(&Album{ArtistID: 2, Name: vars.ReleaseWildwaysTheX}))

	// action
	albums, err := DbMgr.GetAlbums(3)

	// assert
	assert.NoError(t.T(), err)
	assert.Empty(t.T(), albums)
}
