package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestNotifications_CreateAndGet() {
	// action
	now := time.Now().UTC()
	err := Mgr.CreateNotification(&Notification{UserName: vars.UserObjque, Date: now, ReleaseID: 1})

	// assert
	assert.NoError(t.T(), err)
	notifications, err := Mgr.GetNotificationsForUser(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), notifications, 1)
	assert.Equal(t.T(), vars.UserObjque, notifications[0].UserName)
	assert.Equal(t.T(), now, notifications[0].Date)
	assert.Equal(t.T(), uint64(1), notifications[0].ReleaseID)
}
