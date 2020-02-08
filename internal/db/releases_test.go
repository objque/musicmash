package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleases_EnsureExists() {
	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		StoreName: vars.StoreDeezer,
		StoreID:   vars.StoreApple,
		Explicit:  true,
	})

	// assert
	assert.NoError(t.T(), err)
	releases, err := DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.True(t.T(), releases[0].Explicit)
}

func (t *testDBSuite) TestReleases_FindReleases() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  vars.StoreIDW,
		StoreName: vars.StoreApple,
		StoreID:   vars.StoreIDA,
		Poster:    vars.PosterSimple,
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  vars.StoreIDQ,
		StoreName: vars.StoreApple,
		StoreID:   vars.StoreIDB,
		Explicit:  true,
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[0].ArtistID)
	assert.True(t.T(), releases[0].Explicit)
}
