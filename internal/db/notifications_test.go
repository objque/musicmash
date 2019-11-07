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
