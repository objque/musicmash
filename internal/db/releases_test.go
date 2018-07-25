package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Releases_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    1265867295,
	})

	// assert
	assert.NoError(t, err)
	release, err := DbMgr.FindRelease("skrillex", 1265867295)
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", release.ArtistName)
	assert.Equal(t, uint64(1265867295), release.StoreID)
}

func TestDB_Releases_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    uint64(1265867295),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    uint64(63569098),
	}))

	// action
	releases, err := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Equal(t, "S.P.Y", releases[1].ArtistName)
}
