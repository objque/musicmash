package db

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Notifications_CreateAndGet(t *testing.T) {
	setup()
	defer teardown()

	// action
	now := time.Now().UTC()
	err := DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 1})

	// assert
	assert.NoError(t, err)
	notifications, err := DbMgr.GetNotificationsForUser(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, testutil.UserObjque, notifications[0].UserName)
	assert.Equal(t, now, notifications[0].Date)
	assert.Equal(t, uint64(1), notifications[0].ReleaseID)
}

func TestDB_Notifications_CreateAndGetWithID(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 1}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserBot, Date: now, ReleaseID: 1}))

	// action
	notification, err := DbMgr.IsUserNotified(testutil.UserObjque, 1, false)

	// action
	assert.NoError(t, err)
	assert.Equal(t, testutil.UserObjque, notification.UserName)
	assert.Equal(t, now, notification.Date)
	assert.Equal(t, uint64(1), notification.ReleaseID)
}

func TestDB_Notifications_CreateAndGetWithID_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 1}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 2}))

	// action
	notification, err := DbMgr.IsUserNotified(testutil.UserBot, 1, false)

	// action
	assert.Error(t, err)
	assert.Nil(t, notification)
}

func TestDB_Notifications_AlreadyNotifiedAboutFuture_ButReleasedToday(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	// user was notified about coming release month ago
	assert.NoError(t, DbMgr.CreateNotification(&Notification{
		UserName:  testutil.UserObjque,
		Date:      now.AddDate(0, -1, 0),
		ReleaseID: 1,
		IsComing:  true,
	}))

	// action
	notification, err := DbMgr.IsUserNotified(testutil.UserObjque, 1, false)

	// assert
	assert.Error(t, err)
	assert.Nil(t, notification)
}
