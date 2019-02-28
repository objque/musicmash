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

func TestDB_Releases_GetReleasesForUserFilterByPeriod(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC()
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(testutil.UserObjque, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistSkrillex, releases[0].ArtistName)
}

func TestDB_Releases_GetReleasesForUserFilterByPeriod_WithFuture(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC().Add(time.Hour * 48)
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(testutil.UserObjque, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, testutil.ArtistSkrillex, releases[0].ArtistName)
	assert.Equal(t, testutil.ArtistSPY, releases[1].ArtistName)
}

func TestDB_Releases_GetReleasesForUserSince(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	releases, err := DbMgr.GetReleasesForUserSince(testutil.UserObjque, time.Now())

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistSPY, releases[0].ArtistName)
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

func TestDB_Releases_FindNewReleasesForUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  now.Add(-time.Hour * 48),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  now.Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	releases, err := DbMgr.FindNewReleasesForUser(testutil.UserObjque, now)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, testutil.ArtistSPY, releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 1)
}

func TestDB_Releases_FindNewReleasesForUser_ThatWasAnnouncedEarlier(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{
		testutil.ArtistSkrillex,
		testutil.ArtistSPY,
		testutil.ArtistArchitects,
		testutil.ArtistWildways,
		testutil.ArtistAlgorithm,
	}))
	//
	// shouldn't be in the output
	//
	// the oldest album that was released
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreApple,
		StoreID:    "10000",
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Add(-testutil.Month * 2),
	}))
	// announced album that was found many time ago
	// because it was added in the past and already was notified
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    "10002",
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Add(testutil.Month),
	}))
	// announced album that was found today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistTomOdell,
		StoreName:  testutil.StoreApple,
		StoreID:    "10502",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24).Add(testutil.Month),
	}))
	// album that was found today and released today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistRitaOra,
		StoreName:  testutil.StoreApple,
		StoreID:    "10503",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))
	//
	// must be in the output
	//
	// album that was announced too many time ago and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Truncate(time.Hour * 24),
	}))
	// announced album that was found today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistArchitects,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		CreatedAt:  now,
		Released:   now.Add(testutil.Month),
	}))
	// album that was found today and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistAlgorithm,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))

	// action
	releases, err := DbMgr.FindNewReleasesForUser(testutil.UserObjque, now.Truncate(time.Hour*24))

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 3)
	want := map[string]string{
		testutil.ArtistSkrillex:   testutil.StoreIDA,
		testutil.ArtistArchitects: testutil.StoreIDB,
		testutil.ArtistAlgorithm:  testutil.StoreIDC,
	}
	for _, release := range releases {
		assert.Equal(t, want[release.ArtistName], release.StoreID)
		assert.Len(t, release.Stores, 1)
	}
}

func TestDB_Releases_FindNewReleasesForUser_ExcludeAlreadyDelivered_WithAnotherConditions(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{
		testutil.ArtistSkrillex,
		testutil.ArtistSPY,
		testutil.ArtistArchitects,
		testutil.ArtistWildways,
		testutil.ArtistAlgorithm,
	}))
	//
	// shouldn't be in the output
	//
	// the oldest album that was released
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreApple,
		StoreID:    "10000",
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Add(-testutil.Month * 2),
	}))
	// announced album that was found many time ago
	// because it was added in the past and already was notified
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    "10002",
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Add(testutil.Month),
	}))
	// announced album that was found today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistTomOdell,
		StoreName:  testutil.StoreApple,
		StoreID:    "10502",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24).Add(testutil.Month),
	}))
	// album that was found today and released today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistRitaOra,
		StoreName:  testutil.StoreApple,
		StoreID:    "10503",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))
	// album that was found today and released today
	// but user already got notification
	release := Release{
		ArtistName: testutil.ArtistWildways,
		StoreName:  testutil.StoreApple,
		StoreID:    "10723",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&release))
	DbMgr.MarkReleasesAsDelivered(testutil.UserObjque, []*Release{&release})
	//
	// must be in the output
	//
	// album that was announced too many time ago and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		CreatedAt:  now.Add(-testutil.Month * 3),
		Released:   now.Truncate(time.Hour * 24),
	}))
	// announced album that was found today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistArchitects,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		CreatedAt:  now,
		Released:   now.Add(testutil.Month),
	}))
	// album that was found today and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistAlgorithm,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDC,
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))

	// action
	releases, err := DbMgr.FindNewReleasesForUser(testutil.UserObjque, now.Truncate(time.Hour*24))

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 3)
	want := map[string]string{
		testutil.ArtistSkrillex:   testutil.StoreIDA,
		testutil.ArtistArchitects: testutil.StoreIDB,
		testutil.ArtistAlgorithm:  testutil.StoreIDC,
	}
	for _, release := range releases {
		assert.Equal(t, want[release.ArtistName], release.StoreID)
		assert.Len(t, release.Stores, 1)
	}
}

func TestDB_Releases_FindNewReleasesForUser_ExcludeAlreadyDelivered(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{
		testutil.ArtistArchitects,
		testutil.ArtistAlgorithm,
	}))
	architectsRelease := Release{
		ArtistName: testutil.ArtistArchitects,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		CreatedAt:  now,
		// announced album that was found today
		Released: now.Add(testutil.Month),
	}
	algorithmRelease := Release{
		ArtistName: testutil.ArtistAlgorithm,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		CreatedAt:  now,
		// album that was found today and released today
		Released: now.Truncate(time.Hour * 24),
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&architectsRelease))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&algorithmRelease))
	releases, err := DbMgr.FindNewReleasesForUser(testutil.UserObjque, now.Truncate(time.Hour*24))
	assert.NoError(t, err)
	assert.Len(t, releases, 2)

	// action
	DbMgr.MarkReleasesAsDelivered(testutil.UserObjque, []*Release{&architectsRelease, &algorithmRelease})
	releases, err = DbMgr.FindNewReleasesForUser(testutil.UserObjque, now.Truncate(time.Hour*24))

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 0)
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
