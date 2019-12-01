package db

import (
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleases_EnsureExists() {
	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		StoreName: testutil.StoreDeezer,
		StoreID:   testutil.StoreApple,
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
		ArtistID:  testutil.StoreIDW,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDA,
		Poster:    testutil.PosterSimple,
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDQ,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDB,
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), releases[0].ArtistID)
}
