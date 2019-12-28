package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestArtists_Search() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutils.ArtistArchitects}))

	// action
	releases, err := artists.Search(t.client, "arch")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), testutils.ArtistArchitects, releases[0].Name)
}

func (t *testAPISuite) TestArtists_Search_Empty() {
	// action
	releases, err := artists.Search(t.client, "arch")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}

func (t *testAPISuite) TestArtists_Search_Internal() {
	// arrange
	// close connection manually to get internal error
	assert.NoError(t.T(), db.DbMgr.Close())

	// action
	releases, err := artists.Search(t.client, "arch")

	// assert
	assert.Error(t.T(), err)
	assert.Len(t.T(), releases, 0)
}

func (t *testAPISuite) TestArtists_Create() {
	// action
	artist := &artists.Artist{Name: testutils.ArtistArchitects}
	err := artists.Create(t.client, artist)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testutils.ArtistArchitects, artist.Name)
	assert.Equal(t.T(), int64(1), artist.ID)
}

func (t *testAPISuite) TestArtists_Associate() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutils.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(testutils.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutils.StoreApple, StoreID: testutils.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testutils.StoreIDA, info.StoreID)
	assert.Equal(t.T(), testutils.StoreApple, info.StoreName)
	assert.Equal(t.T(), int64(1), info.ArtistID)
}

func (t *testAPISuite) TestArtists_Associate_ArtistNotFound() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(testutils.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutils.StoreApple, StoreID: testutils.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Associate_AlreadyAssociated() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutils.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(testutils.StoreApple))
	assert.NoError(t.T(), db.DbMgr.EnsureAssociationExists(1, testutils.StoreApple, testutils.StoreIDA))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutils.StoreApple, StoreID: testutils.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Associate_StoreNotFound() {
	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutils.StoreApple, StoreID: testutils.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Get() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutils.ArtistArchitects}))

	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), testutils.ArtistArchitects, artist.Name)
	assert.Empty(t.T(), artist.Albums)
}

func (t *testAPISuite) TestArtists_GetWithAlbums() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutils.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureAlbumExists(&db.Album{ArtistID: 1, Name: testutils.ReleaseAlgorithmFloatingIP}))

	// action
	artist, err := artists.Get(t.client, 1, &artists.GetOptions{WithAlbums: true})

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), testutils.ArtistArchitects, artist.Name)
	assert.Len(t.T(), artist.Albums, 1)
	assert.Equal(t.T(), testutils.ReleaseAlgorithmFloatingIP, artist.Albums[0].Name)
}

func (t *testAPISuite) TestArtists_Get_NotFound() {
	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.Error(t.T(), err)
	assert.Nil(t.T(), artist)
}
