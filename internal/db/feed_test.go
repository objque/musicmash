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
	const userID = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	feed, err := DbMgr.GetUserFeedSince(userID, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, feed.Announced, 1)
	assert.Len(t, feed.Released, 1)
	assert.Equal(t, "S.P.Y", feed.Announced[0].ArtistName)
	assert.Equal(t, "skrillex", feed.Released[0].ArtistName)
}
