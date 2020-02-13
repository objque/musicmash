package api

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestArtists_Search() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))

	// action
	artists, err := artists.Search(t.client, "arch")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
	assert.Equal(t.T(), vars.ArtistArchitects, artists[0].Name)
}

func (t *testAPISuite) TestArtists_Search_Empty() {
	// action
	artists, err := artists.Search(t.client, "arch")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 0)
}

func (t *testAPISuite) TestArtists_Search_Internal() {
	// arrange
	// close connection manually to get internal error
	assert.NoError(t.T(), db.DbMgr.Close())

	// action
	artists, err := artists.Search(t.client, "arch")

	// assert
	assert.Error(t.T(), err)
	assert.Len(t.T(), artists, 0)
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

func (t *testAPISuite) TestArtists_Get_Releases() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: 666}))
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       vars.StoreIDQ,
		ArtistID: 666,
		Title:    vars.ArtistArchitects,
		Released: time.Now(),
	}))

	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       vars.StoreIDW,
		ArtistID: 777,
		Title:    vars.ArtistArchitects,
		Released: time.Now().UTC().AddDate(-1, 0, 0),
	}))

	// action
	releases, err := releases.By(t.client, 666)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), vars.ArtistArchitects, releases[0].Title)
	assert.Equal(t.T(), uint64(vars.StoreIDQ), releases[0].ID)
}

func (t *testAPISuite) TestArtists_Get_Releases_Empty() {
	// action
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: 666}))
	releases, err := releases.By(t.client, 666)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}

func (t *testAPISuite) TestArtists_Get_Releases_Internal() {
	// arrange
	_ = db.DbMgr.Close()

	// action
	releases, err := releases.By(t.client, 666)

	// assert
	assert.Error(t.T(), err)
	assert.Len(t.T(), releases, 0)
}

func (t *testAPISuite) TestArtists_Get_Releases_ArtistNotFound() {
	// action
	releases, err := releases.By(t.client, 666)

	// assert
	assert.Error(t.T(), err)
	assert.Len(t.T(), releases, 0)
}
