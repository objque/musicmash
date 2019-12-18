package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestInternalNotifications_Find() {
	const (
		notificationService = "email"
		notificationData    = "duplicated-mail@inbox.me"
	)

	// arrange
	// create artist
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ, Name: testutil.ArtistArchitects}))
	// subscribe users
	assert.NoError(t.T(), DbMgr.CreateSubscription(&Subscription{
		CreatedAt: time.Now().UTC().AddDate(-1, 0, 0),
		UserName:  testutil.UserObjque,
		ArtistID:  testutil.StoreIDQ,
	}))
	assert.NoError(t.T(), DbMgr.CreateSubscription(&Subscription{
		CreatedAt: time.Now().UTC().AddDate(-1, 0, 0),
		UserName:  testutil.UserBot,
		ArtistID:  testutil.StoreIDQ,
	}))
	// fill contacts
	assert.NoError(t.T(), DbMgr.EnsureNotificationServiceExists(notificationService))
	assert.NoError(t.T(), DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserObjque, Service: notificationService, Data: notificationData,
	}))
	assert.NoError(t.T(), DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserBot, Service: notificationService, Data: notificationData,
	}))
	// fill releases
	// won't be in output
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        100,
		ArtistID:  testutil.StoreIDQ,
		Released:  time.Now().UTC().AddDate(0, -1, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        105,
		ArtistID:  testutil.StoreIDQ,
		Released:  time.Now().UTC().AddDate(0, 10, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-future-release-wont-be-in-output",
	}))
	// deliver notifications
	assert.NoError(t.T(), DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: true, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: false, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), DbMgr.CreateNotification(&Notification{
		ReleaseID: 100, IsComing: false, UserName: testutil.UserBot, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), DbMgr.CreateNotification(&Notification{
		ReleaseID: 105, IsComing: true, UserName: testutil.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), DbMgr.CreateNotification(&Notification{
		ReleaseID: 105, IsComing: true, UserName: testutil.UserBot, Date: time.Now().UTC(),
	}))
	// fill releases
	// should be in output
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        20,
		ArtistID:  testutil.StoreIDQ,
		Title:     testutil.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: testutil.StoreApple,
		StoreID:   "this-oldest-release-have-to-be-in-output",
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
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
	assert.NoError(t.T(), err)
	// 2 oldest releases weren't delivery
	// 2 coming releases weren't delivery
	assert.Len(t.T(), notifications, 4)
	for _, notification := range notifications {
		assert.Equal(t.T(), int64(testutil.StoreIDQ), notification.ArtistID)
		assert.Equal(t.T(), testutil.ArtistArchitects, notification.ArtistName)
		assert.Equal(t.T(), notificationService, notification.Service)
		assert.Equal(t.T(), notificationData, notification.Data)
		assert.Equal(t.T(), testutil.StoreApple, notification.StoreName)
		assert.Contains(t.T(), notification.StoreID, "have-to-be-in-output")
	}
}

func (t *testDBSuite) TestInternalNotifications_SubscribedAfterRelease() {
	// user should not receive notifications about releases
	// which released before he subscribed for an artist

	// arrange
	// create artist
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ, Name: testutil.ArtistArchitects}))
	// subscribe users
	assert.NoError(t.T(), DbMgr.CreateSubscription(&Subscription{
		// by default created_at is now()
		UserName: testutil.UserObjque,
		ArtistID: testutil.StoreIDQ,
	}))
	// fill releases
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDQ,
		Title:     testutil.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: testutil.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
	}))
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(&Release{
		ArtistID:  testutil.StoreIDQ,
		Title:     testutil.ArtistAlgorithm,
		Released:  time.Now().UTC().AddDate(1, 0, 0),
		StoreName: testutil.StoreApple,
		StoreID:   "this-future-release-have-to-be-in-output",
		Poster:    testutil.PosterSimple,
	}))

	// action
	notifications, err := DbMgr.FindNotReceivedNotifications()

	// assert
	assert.NoError(t.T(), err)
	// user subscribed now(), so:
	// should receive notification about future release
	// shouldn't receive notification about old release
	assert.Len(t.T(), notifications, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), notifications[0].ArtistID)
	assert.Equal(t.T(), testutil.ArtistArchitects, notifications[0].ArtistName)
	assert.Equal(t.T(), testutil.ArtistAlgorithm, notifications[0].Title)
	assert.Equal(t.T(), testutil.StoreApple, notifications[0].StoreName)
	assert.Equal(t.T(), testutil.PosterSimple, notifications[0].ReleasePoster)
	assert.Contains(t.T(), notifications[0].StoreID, "this-future-release-have-to-be-in-output")
}
