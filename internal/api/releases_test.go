package api

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestReleases_Get() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       vars.StoreIDQ,
		Title:    vars.ArtistArchitects,
		Released: time.Now(),
	}))

	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       vars.StoreIDW,
		Title:    vars.ArtistArchitects,
		Released: time.Now().UTC().AddDate(-1, 0, 0),
	}))

	// action
	since := time.Now().UTC().AddDate(0, -1, 0)
	releases, err := releases.Get(t.client, since)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), vars.ArtistArchitects, releases[0].Title)
	assert.Equal(t.T(), uint64(vars.StoreIDQ), releases[0].ID)
}

func (t *testAPISuite) TestReleases_Get_Empty() {
	// action
	since := time.Now().UTC().AddDate(0, -1, 0)
	releases, err := releases.Get(t.client, since)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}

func (t *testAPISuite) TestReleases_Get_Internal() {
	// arrange
	_ = db.DbMgr.Close()

	// action
	since := time.Now().UTC().AddDate(0, -1, 0)
	releases, err := releases.Get(t.client, since)

	// assert
	assert.Error(t.T(), err)
	assert.Len(t.T(), releases, 0)
}
