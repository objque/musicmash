package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_Releases_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		StoreName: "deezer",
		StoreID:   "xyz",
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
	const userName = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(userName, []string{"skrillex", "S.P.Y"}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC()
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(userName, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
}

func TestDB_Releases_GetReleasesForUserFilterByPeriod_WithFuture(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userName = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(userName, []string{"skrillex", "S.P.Y"}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC().Add(time.Hour * 48)
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(userName, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
}

func TestDB_Releases_GetReleasesForUserSince(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userName = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(userName, []string{"skrillex", "S.P.Y"}))

	// action
	releases, err := DbMgr.GetReleasesForUserSince(userName, time.Now())

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
}

func TestDB_Releases_FindReleasesWithFilter(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  time.Now().UTC().Add(-time.Hour * 48),
	}))
	date := time.Now().UTC().Add(-time.Hour * 24)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  date,
	}))

	// action
	releases, err := DbMgr.FindNewReleases(date)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
}

func TestDB_Releases_FindNewReleasesForUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  now.Add(-time.Hour * 48),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  now.Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists("objque", []string{"skrillex", "S.P.Y"}))

	// action
	releases, err := DbMgr.FindNewReleasesForUser("objque", now)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 1)
}

func TestDB_Releases_FindNewReleasesForUser_ThatWasAnnouncedEarlier(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const (
		day      = time.Hour * 24
		month    = day * 30
		userName = "objque@me"
	)
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.SubscribeUserForArtists(userName, []string{
		"Skrillex",
		"S.P.Y",
		"Architects",
		"Wildways",
		"The Algorithm",
	}))
	//
	// shouldn't be in the output
	//
	// the oldest album that was released
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "Wildways",
		StoreName:  "itunes",
		StoreID:    "10000",
		CreatedAt:  now.Add(-month * 3),
		Released:   now.Add(-month * 2),
	}))
	// announced album that was found many time ago
	// because it was added in the past and already was notified
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "10002",
		CreatedAt:  now.Add(-month * 3),
		Released:   now.Add(month),
	}))
	// announced album that was found today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "Tom Odell",
		StoreName:  "itunes",
		StoreID:    "10502",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24).Add(month),
	}))
	// album that was found today and released today
	// but user not subscribed for this artist
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "RITA ORA",
		StoreName:  "itunes",
		StoreID:    "10503",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))
	//
	// must be in the output
	//
	// album that was announced too many time ago and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "Skrillex",
		StoreName:  "itunes",
		StoreID:    "30001",
		CreatedAt:  now.Add(-month * 3),
		Released:   now.Truncate(time.Hour * 24),
	}))
	// announced album that was found today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "Architects",
		StoreName:  "itunes",
		StoreID:    "30002",
		CreatedAt:  now,
		Released:   now.Add(month),
	}))
	// album that was found today and released today
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "The Algorithm",
		StoreName:  "itunes",
		StoreID:    "30003",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}))

	// action
	releases, err := DbMgr.FindNewReleasesForUser(userName, now.Truncate(time.Hour*24))

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 3)
	want := map[string]string{"Skrillex": "30001", "Architects": "30002", "The Algorithm": "30003"}
	for _, release := range releases {
		assert.Equal(t, want[release.ArtistName], release.StoreID)
		assert.Len(t, release.Stores, 1)
	}
}

func TestDB_Releases_FindReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Poster:     "http://pic.jpeg",
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
	}))

	// action
	releases, err := DbMgr.FindReleases(map[string]interface{}{"poster": ""})

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
}

func TestDB_Releases_UpdateRelease(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	release := Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&release))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
	}))

	// action
	release.Poster = "http://pic.jpeg"
	err := DbMgr.UpdateRelease(&release)
	releases, _ := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Equal(t, "itunes", releases[0].StoreName)
	assert.Equal(t, "182821355", releases[0].StoreID)
	assert.Equal(t, "http://pic.jpeg", releases[0].Poster)
	// another release must not change
	assert.Equal(t, "S.P.Y", releases[1].ArtistName)
	assert.Equal(t, "itunes", releases[1].StoreName)
	assert.Equal(t, "213551828", releases[1].StoreID)
}
