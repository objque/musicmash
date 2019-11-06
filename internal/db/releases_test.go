package db

import (
	"testing"
	"time"

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

func TestDB_Releases_FindArtistsWithNewReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ArtistID:  testutil.StoreIDW,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDA,
		Poster:    testutil.PosterSimple,
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC().AddDate(0, -1, 0),
		ArtistID:  testutil.StoreIDQ,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDB,
	}))

	// action
	artists, err := DbMgr.FindArtistsWithNewReleases(time.Now().UTC().Add(-time.Hour))

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
	assert.Equal(t, int64(testutil.StoreIDW), artists[0])
}

func TestDB_Releases_UpdateRelease(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	release := Release{
		ArtistID:  testutil.StoreIDW,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDA,
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&release))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDQ,
		StoreName: testutil.StoreApple,
		StoreID:   testutil.StoreIDB,
	}))

	// action
	release.Poster = testutil.PosterSimple
	err := DbMgr.UpdateRelease(&release)
	releases, _ := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, int64(testutil.StoreIDW), releases[0].ArtistID)
	assert.Equal(t, testutil.StoreApple, releases[0].StoreName)
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, testutil.PosterSimple, releases[0].Poster)
	// another release must not change
	assert.Equal(t, int64(testutil.StoreIDQ), releases[1].ArtistID)
	assert.Equal(t, testutil.StoreApple, releases[1].StoreName)
	assert.Equal(t, testutil.StoreIDB, releases[1].StoreID)
}
