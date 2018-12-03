package db

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Notifications_MarkAndGet(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	release := Release{
		ArtistName: testutil.ArtistArchitects,
		StoreName:  testutil.StoreApple,
		StoreID:    "30002",
		CreatedAt:  now,
		Released:   now.Truncate(time.Hour * 24),
	}
	assert.NoError(t, DbMgr.EnsureReleaseExists(&release))

	// action
	DbMgr.MarkReleasesAsDelivered(testutil.UserObjque, []*Release{&release})

	// assert
	notifications, err := DbMgr.GetNotificationsForUser(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
}
