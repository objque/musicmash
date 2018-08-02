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
		StoreID:    182821355,
	})

	// assert
	assert.NoError(t, err)
	release, err := DbMgr.FindRelease("skrillex", 182821355)
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", release.ArtistName)
	assert.Equal(t, uint64(182821355), release.StoreID)
}

func TestDB_Releases_IsExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	const storeID = uint64(182821355)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    storeID,
	}))

	// action and assert
	assert.True(t, DbMgr.IsReleaseExists(storeID))
}

func TestDB_Releases_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
	}))

	// action
	releases, err := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Equal(t, "S.P.Y", releases[1].ArtistName)
}
