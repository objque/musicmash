package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleases_EnsureExists() {
	// action
	err := Mgr.EnsureReleaseExists(&Release{
		StoreName: vars.StoreDeezer,
		StoreID:   vars.StoreApple,
		Explicit:  true,
	})

	// assert
	assert.NoError(t.T(), err)
	releases, err := Mgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.True(t.T(), releases[0].Explicit)
}

func (t *testDBSuite) TestReleases_FindReleases() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  vars.StoreIDW,
		StoreName: vars.StoreApple,
		StoreID:   vars.StoreIDA,
		Title:     vars.ArtistAlgorithm,
		Explicit:  true,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  vars.StoreIDQ,
		StoreName: vars.StoreApple,
		StoreID:   vars.StoreIDB,
	}))

	// action
	releases, err := Mgr.FindReleases(vars.StoreIDW, vars.ArtistAlgorithm)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(vars.StoreIDW), releases[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistAlgorithm, releases[0].Title)
	assert.True(t.T(), releases[0].Explicit)
}
