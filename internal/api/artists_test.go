package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestArtists_Search() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))

	// action
	releases, err := artists.Search(t.client, "arch")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), vars.ArtistArchitects, releases[0].Name)
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
	artist := &artists.Artist{Name: vars.ArtistArchitects}
	err := artists.Create(t.client, artist)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.ArtistArchitects, artist.Name)
	assert.Equal(t.T(), int64(1), artist.ID)
}

func (t *testAPISuite) TestArtists_Associate() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.StoreIDA, info.StoreID)
	assert.Equal(t.T(), vars.StoreApple, info.StoreName)
	assert.Equal(t.T(), int64(1), info.ArtistID)
}

func (t *testAPISuite) TestArtists_Associate_ArtistNotFound() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Associate_AlreadyAssociated() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), db.DbMgr.EnsureAssociationExists(1, vars.StoreApple, vars.StoreIDA))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Associate_StoreNotFound() {
	// action
	info := &artists.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := artists.Associate(t.client, info)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestArtists_Get() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))

	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), vars.ArtistArchitects, artist.Name)
}

func (t *testAPISuite) TestArtists_Get_NotFound() {
	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.Error(t.T(), err)
	assert.Nil(t.T(), artist)
}
