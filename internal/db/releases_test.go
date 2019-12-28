package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleases_EnsureExists() {
	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		StoreName: testutils.StoreDeezer,
		StoreID:   testutils.StoreApple,
	})

	// assert
	assert.NoError(t.T(), err)
	releases, err := DbMgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
}

func (t *testDBSuite) TestReleases_FindReleases() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutils.StoreIDW,
		StoreName: testutils.StoreApple,
		StoreID:   testutils.StoreIDA,
		Poster:    testutils.PosterSimple,
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutils.StoreIDQ,
		StoreName: testutils.StoreApple,
		StoreID:   testutils.StoreIDB,
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), releases[0].ArtistID)
}
