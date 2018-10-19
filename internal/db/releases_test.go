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
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "S.P.Y",
	}))

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
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "S.P.Y",
	}))

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
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserName:   userName,
		ArtistName: "S.P.Y",
	}))

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
	const userID = "objque@me"
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
	const userID = "objque@me"
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
}
