package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_NotificationSettings_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "email",
		Data:     "email@test.io",
	})

	// assert
	assert.NoError(t, err)
	settings, err := DbMgr.FindNotificationSettings(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, settings, 1)
	assert.Equal(t, testutil.UserObjque, settings[0].UserName)
	assert.Equal(t, "email", settings[0].Service)
	assert.Equal(t, "email@test.io", settings[0].Data)
}

func TestDB_NotificationSettings_Update(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "email",
		Data:     "email@test.io",
	}))
	assert.NoError(t, DbMgr.EnsureNotificationSettingsExists(&NotificationSettings{
		UserName: testutil.UserBot,
		Service:  "email",
		Data:     "email@test.io",
	}))

	// action
	err := DbMgr.UpdateNotificationSettings(&NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "email",
		Data:     "objque@test.io",
	})

	// assert
	assert.NoError(t, err)
	settings, err := DbMgr.FindNotificationSettings(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, settings, 1)
	assert.Equal(t, testutil.UserObjque, settings[0].UserName)
	assert.Equal(t, "email", settings[0].Service)
	assert.Equal(t, "objque@test.io", settings[0].Data)
}
