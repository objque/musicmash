package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_Feed_GetUserFeedSince(t *testing.T) {
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
	feed, err := DbMgr.GetUserFeedSince(userName, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, feed.Announced, 1)
	assert.Len(t, feed.Released, 1)
	assert.Equal(t, "S.P.Y", feed.Announced[0].ArtistName)
	assert.Equal(t, "skrillex", feed.Released[0].ArtistName)
}
