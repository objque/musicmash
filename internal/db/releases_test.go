package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Releases_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		StoreName: testutil.StoreDeezer,
		StoreID:   testutil.StoreApple,
	})

	// assert
	assert.NoError(t, err)
	releases, err := DbMgr.GetAllReleases()
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
}

func TestDB_Releases_FindReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDW,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDA,
		Poster:    testutil.PosterSimple,
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDQ,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDB,
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), releases[0].ArtistID)
}
