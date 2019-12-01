package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestNotifications_CreateAndGet() {
	// action
	now := time.Now().UTC()
	err := DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 1})

	// assert
	assert.NoError(t.T(), err)
	notifications, err := DbMgr.GetNotificationsForUser(testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), notifications, 1)
	assert.Equal(t.T(), testutil.UserObjque, notifications[0].UserName)
	assert.Equal(t.T(), now, notifications[0].Date)
	assert.Equal(t.T(), uint64(1), notifications[0].ReleaseID)
}
