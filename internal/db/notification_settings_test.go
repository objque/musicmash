package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestNotificationSettings_EnsureExists() {
	// action
	err := DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutils.UserObjque,
		Service:  "email",
		Data:     "email@test.io",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := DbMgr.FindNotificationSettings(testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 1)
	assert.Equal(t.T(), testutils.UserObjque, settings[0].UserName)
	assert.Equal(t.T(), "email", settings[0].Service)
	assert.Equal(t.T(), "email@test.io", settings[0].Data)
}

func (t *testDBSuite) TestNotificationSettings_Update() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutils.UserObjque,
		Service:  "email",
		Data:     "email@test.io",
	}))
	assert.NoError(t.T(), DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutils.UserBot,
		Service:  "email",
		Data:     "email@test.io",
	}))

	// action
	err := DbMgr.UpdateNotificationSettings(&NotificationSettings{
		UserName: testutils.UserObjque,
		Service:  "email",
		Data:     "objque@test.io",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := DbMgr.FindNotificationSettings(testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 1)
	assert.Equal(t.T(), testutils.UserObjque, settings[0].UserName)
	assert.Equal(t.T(), "email", settings[0].Service)
	assert.Equal(t.T(), "objque@test.io", settings[0].Data)
}
