package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestInternalNotifications_Find() {
	const (
		notificationService = "email"
		notificationData    = "duplicated-mail@inbox.me"
	)

	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	// create artist
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 1, Name: vars.ArtistArchitects}))
	// subscribe users
	assert.NoError(t.T(), Mgr.CreateSubscription(&Subscription{
		CreatedAt: time.Now().UTC().AddDate(-1, 0, 0),
		UserName:  vars.UserObjque,
		ArtistID:  1,
	}))
	assert.NoError(t.T(), Mgr.CreateSubscription(&Subscription{
		CreatedAt: time.Now().UTC().AddDate(-1, 0, 0),
		UserName:  vars.UserBot,
		ArtistID:  1,
	}))
	// fill contacts
	assert.NoError(t.T(), Mgr.EnsureNotificationServiceExists(notificationService))
	assert.NoError(t.T(), Mgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: vars.UserObjque, Service: notificationService, Data: notificationData,
	}))
	assert.NoError(t.T(), Mgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: vars.UserBot, Service: notificationService, Data: notificationData,
	}))
	// fill releases
	// won't be in output
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        1,
		ArtistID:  1,
		Released:  time.Now().UTC().AddDate(0, -1, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        2,
		ArtistID:  1,
		Released:  time.Now().UTC().AddDate(0, 10, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-future-release-wont-be-in-output",
	}))
	// deliver notifications
	assert.NoError(t.T(), Mgr.CreateNotification(&Notification{
		ReleaseID: 1, IsComing: true, UserName: vars.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), Mgr.CreateNotification(&Notification{
		ReleaseID: 1, IsComing: false, UserName: vars.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), Mgr.CreateNotification(&Notification{
		ReleaseID: 1, IsComing: false, UserName: vars.UserBot, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), Mgr.CreateNotification(&Notification{
		ReleaseID: 2, IsComing: true, UserName: vars.UserObjque, Date: time.Now().UTC(),
	}))
	assert.NoError(t.T(), Mgr.CreateNotification(&Notification{
		ReleaseID: 2, IsComing: true, UserName: vars.UserBot, Date: time.Now().UTC(),
	}))
	// fill releases
	// should be in output
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        3,
		ArtistID:  1,
		Title:     vars.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: vars.StoreApple,
		Type:      vars.ReleaseTypeAlbum,
		Explicit:  true,
		StoreID:   "this-oldest-release-have-to-be-in-output",
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC(),
		ID:        4,
		ArtistID:  1,
		Title:     vars.ArtistAlgorithm,
		Released:  time.Now().UTC().AddDate(1, 0, 0),
		StoreName: vars.StoreApple,
		Type:      vars.ReleaseTypeAlbum,
		Explicit:  true,
		StoreID:   "this-future-release-have-to-be-in-output",
	}))

	// action
	notifications, err := Mgr.FindNotReceivedNotifications()

	// assert
	assert.NoError(t.T(), err)
	// 2 oldest releases weren't delivery
	// 2 coming releases weren't delivery
	assert.Len(t.T(), notifications, 4)
	for _, notification := range notifications {
		assert.Equal(t.T(), int64(1), notification.ArtistID)
		assert.Equal(t.T(), vars.ArtistArchitects, notification.ArtistName)
		assert.Equal(t.T(), notificationService, *notification.Service)
		assert.Equal(t.T(), notificationData, *notification.Data)
		assert.Equal(t.T(), vars.StoreApple, notification.StoreName)
		assert.Equal(t.T(), vars.ReleaseTypeAlbum, notification.Type)
		assert.Contains(t.T(), notification.StoreID, "have-to-be-in-output")
		assert.True(t.T(), notification.Explicit)
	}
}

func (t *testDBSuite) TestInternalNotifications_SubscribedAfterRelease() {
	// user should not receive notifications about releases
	// which released before he subscribed for an artist

	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	// create artist
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects}))
	// subscribe users
	assert.NoError(t.T(), Mgr.CreateSubscription(&Subscription{
		CreatedAt: time.Now().UTC(),
		UserName:  vars.UserObjque,
		ArtistID:  1,
	}))
	// fill releases
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  1,
		Title:     vars.ReleaseArchitectsHollyHell,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: vars.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
		Explicit:  true,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  1,
		Title:     vars.ArtistAlgorithm,
		Released:  time.Now().UTC().AddDate(1, 0, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-future-release-have-to-be-in-output",
		Poster:    vars.PosterSimple,
		Explicit:  true,
	}))

	// action
	notifications, err := Mgr.FindNotReceivedNotifications()

	// assert
	assert.NoError(t.T(), err)
	// user subscribed now(), so:
	// should receive notification about future release
	// shouldn't receive notification about old release
	assert.Len(t.T(), notifications, 1)
	assert.Equal(t.T(), int64(1), notifications[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistArchitects, notifications[0].ArtistName)
	assert.Equal(t.T(), vars.ArtistAlgorithm, notifications[0].Title)
	assert.Nil(t.T(), notifications[0].Service)
	assert.Nil(t.T(), notifications[0].Data)
	assert.Equal(t.T(), vars.StoreApple, notifications[0].StoreName)
	assert.Equal(t.T(), vars.PosterSimple, notifications[0].Poster)
	assert.Contains(t.T(), notifications[0].StoreID, "this-future-release-have-to-be-in-output")
}
