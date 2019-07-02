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

func TestDB_Releases_FindReleasesWithFilter(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  time.Now().UTC().Add(-time.Hour * 48),
	}))
	date := time.Now().UTC().Add(-time.Hour * 24)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  date,
	}))

	// action
	releases, err := DbMgr.FindNewReleases(date)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistSPY, releases[0].ArtistName)
}

func TestDB_Releases_FindReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Poster:     testutil.PosterSimple,
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistSPY, releases[0].ArtistName)
}

func TestDB_Releases_UpdateRelease(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	release := Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&release))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
	}))

	// action
	release.Poster = testutil.PosterSimple
	err := DbMgr.UpdateRelease(&release)
	releases, _ := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testutil.ArtistSkrillex, releases[0].ArtistName)
	assert.Equal(t, testutil.StoreApple, releases[0].StoreName)
	assert.Equal(t, testutil.StoreIDA, releases[0].StoreID)
	assert.Equal(t, testutil.PosterSimple, releases[0].Poster)
	// another release must not change
	assert.Equal(t, testutil.ArtistSPY, releases[1].ArtistName)
	assert.Equal(t, testutil.StoreApple, releases[1].StoreName)
	assert.Equal(t, testutil.StoreIDB, releases[1].StoreID)
}

func TestDB_Releases_FindArtistRecentReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// recent releases
	now := time.Now().UTC().Truncate(testutil.Day)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   now.Add(-testutil.Day),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week),
	}))
	// another artist (should not be in the result)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week * 2),
	}))
	// announced releases
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Day),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Month),
	}))

	// action
	releases, err := DbMgr.FindArtistRecentReleases(testutil.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, testutil.ArtistSkrillex, releases[0].ArtistName)
	assert.Equal(t, testutil.StoreApple, releases[0].StoreName)
	assert.Equal(t, testutil.ArtistSkrillex, releases[1].ArtistName)
	assert.Equal(t, testutil.StoreDeezer, releases[1].StoreName)
}

func TestDB_Releases_FindArtistAnnouncedReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// recent releases
	now := time.Now().UTC().Truncate(testutil.Day)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   now.Add(-testutil.Day),
	}))
	// another artists (should not be in the result)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDB,
		Released:   now.Add(-testutil.Week * 2),
	}))
	// announced releases
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Day),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreDeezer,
		StoreID:    testutil.StoreIDC,
		Released:   now.Add(testutil.Month),
	}))

	// action
	releases, err := DbMgr.FindArtistAnnouncedReleases(testutil.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, testutil.ArtistSkrillex, releases[0].ArtistName)
	assert.Equal(t, testutil.StoreApple, releases[0].StoreName)
	assert.Equal(t, testutil.ArtistSkrillex, releases[1].ArtistName)
	assert.Equal(t, testutil.StoreDeezer, releases[1].StoreName)
}
