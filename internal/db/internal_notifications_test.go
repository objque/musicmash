package db

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_InternalNotifications_Find(t *testing.T) {
	setup()
	defer teardown()

	const (
		notificationService = "email"
		notificationData    = "duplicated-mail@inbox.me"
	)

	// arrange
	// create artist
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ, Name: testutil.ArtistArchitects}))
	// subscribe users
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserBot, []int64{testutil.StoreIDQ}))
	// fill contacts
	assert.NoError(t, DbMgr.EnsureNotificationServiceExists(notificationService))
	assert.NoError(t, DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserObjque, Service: notificationService, Data: notificationData,
	}))
	assert.NoError(t, DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserBot, Service: notificationService, Data: notificationData,
	}))
	// fill releases
	// won't be in output
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        100,
		ArtistID:  testutil.StoreIDQ,
		//Released:  time.Now().UTC().AddDate(0, -1, 0),
		Released:  time.Now().UTC().AddDate(0, -1, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        105,
		ArtistID:  testutil.StoreIDQ,
		//Released:  time.Now().UTC().AddDate(0, 10, 0),
		Released:  time.Now().UTC().AddDate(0, 10, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-future-release-wont-be-in-output",
	}))
	// deliver notifications
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: true, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: false, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: false, UserName: testutil.UserBot, Date: time.Now().UTC(),
	}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		ReleaseID: 105, IsComing: true, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		ReleaseID: 105, IsComing: true, UserName: testutil.UserBot, Date: time.Now().UTC(),
	}))
	// fill releases
	// should be in output
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        20,
		ArtistID:  testutil.StoreIDQ,
		Title:     testutil.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: testutil.StoreApple,
		StoreID:   "this-oldest-release-have-to-be-in-output",
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        25,
		ArtistID:  testutil.StoreIDQ,
		Title:     testutil.ArtistAlgorithm,
		Released:  time.Now().UTC().AddDate(1, 0, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-future-release-have-to-be-in-output",
	}))

	// action
	notifications, err := DbMgr.FindNotReceivedNotifications()

	// assert
	assert.NoError(t, err)
	// 2 oldest releases weren't delivery
	// 2 coming releases weren't delivery
	assert.Len(t, notifications, 4)
	for _, notification := range notifications {
		assert.Equal(t, int64(testutil.StoreIDQ), notification.ArtistID)
		assert.Equal(t, testutil.ArtistArchitects, notification.Name)
		assert.Equal(t, notificationService, notification.Service)
		assert.Equal(t, notificationData, notification.Data)
		assert.Equal(t, testutil.StoreApple, notification.StoreName)
		assert.Contains(t, notification.StoreID, "have-to-be-in-output")
	}
}
