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
		Title:      "scarry monsters",
	})

	// assert
	assert.NoError(t, err)
	release, err := DbMgr.FindRelease("skrillex", "scarry monsters")
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", release.ArtistName)
	assert.Equal(t, "scarry monsters", release.Title)
}

func TestDB_Releases_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		Title:      "scarry monsters",
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		Title:      "We Are 21",
	}))

	// action
	releases, err := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Equal(t, "S.P.Y", releases[1].ArtistName)
}
