package api

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Artists_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutil.ArtistArchitects}))

	// action
	releases, err := artists.Search(client, "arch")

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistArchitects, releases[0].Name)
}

func TestAPI_Artists_Search_Empty(t *testing.T) {
	setup()
	defer teardown()

	// action
	releases, err := artists.Search(client, "arch")

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 0)
}

func TestAPI_Artists_Search_Internal(t *testing.T) {
	setup()
	_ = db.DbMgr.Close()
	defer func() { server.Close() }()

	// action
	releases, err := artists.Search(client, "arch")

	// assert
	assert.Error(t, err)
	assert.Len(t, releases, 0)
}

func TestAPI_Artists_Create(t *testing.T) {
	setup()
	defer teardown()

	// action
	artist := &artists.Artist{Name: testutil.ArtistArchitects}
	err := artists.Create(client, artist)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutil.ArtistArchitects, artist.Name)
	assert.Equal(t, int64(1), artist.ID)
}

func TestAPI_Artists_Associate(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutil.ArtistArchitects}))
	assert.NoError(t, db.DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutil.StoreApple, StoreID: testutil.StoreIDA}
	err := artists.Associate(client, info)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutil.StoreIDA, info.StoreID)
	assert.Equal(t, testutil.StoreApple, info.StoreName)
	assert.Equal(t, int64(1), info.ArtistID)
}

func TestAPI_Artists_Associate_ArtistNotFound(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutil.StoreApple, StoreID: testutil.StoreIDA}
	err := artists.Associate(client, info)

	// assert
	assert.Error(t, err)
}

func TestAPI_Artists_Associate_AlreadyAssociated(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutil.ArtistArchitects}))
	assert.NoError(t, db.DbMgr.EnsureStoreExists(testutil.StoreApple))
	assert.NoError(t, db.DbMgr.EnsureAssociationExists(1, testutil.StoreApple, testutil.StoreIDA))

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutil.StoreApple, StoreID: testutil.StoreIDA}
	err := artists.Associate(client, info)

	// assert
	assert.Error(t, err)
}

func TestAPI_Artists_Associate_StoreNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	info := &artists.Association{ArtistID: 1, StoreName: testutil.StoreApple, StoreID: testutil.StoreIDA}
	err := artists.Associate(client, info)

	// assert
	assert.Error(t, err)
}

func TestAPI_Artists_GetByID(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutil.ArtistArchitects}))

	// action
	artist, err := artists.Get(client, 1, nil)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), artist.ID)
	assert.Equal(t, testutil.ArtistArchitects, artist.Name)
	assert.Empty(t, artist.Albums)
}

func TestAPI_Artists_GetByIDWithAlbums(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: testutil.ArtistArchitects}))
	assert.NoError(t, db.DbMgr.EnsureAlbumExists(&db.Album{ArtistID: 1, Name: testutil.ReleaseAlgorithmFloatingIP}))

	// action
	artist, err := artists.Get(client, 1, &artists.GetOptions{WithAlbums: true})

	// assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), artist.ID)
	assert.Equal(t, testutil.ArtistArchitects, artist.Name)
	assert.Len(t, artist.Albums, 1)
	assert.Equal(t, testutil.ReleaseAlgorithmFloatingIP, artist.Albums[0].Name)
}

func TestAPI_Artists_GetByID_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	artist, err := artists.Get(client, 1, nil)

	// assert
	assert.Error(t, err)
	assert.Nil(t, artist)
}
