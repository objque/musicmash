package api

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Releases_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       testutil.StoreIDQ,
		Title:    testutil.ArtistArchitects,
		Released: time.Now(),
	}))

	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ID:       testutil.StoreIDW,
		Title:    testutil.ArtistArchitects,
		Released: time.Now().UTC().AddDate(-1, 0, 0),
	}))

	// action
	since := time.Now().UTC().AddDate(0, -1, 0)
	releases, err := releases.Get(client, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistArchitects, releases[0].Title)
	assert.Equal(t, uint64(testutil.StoreIDQ), releases[0].ID)
}

func TestAPI_Releases_Get_Empty(t *testing.T) {
	setup()
	defer teardown()

	// action
	since := time.Now().UTC().AddDate(0, -1, 0)
	releases, err := releases.Get(client, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 0)
}
