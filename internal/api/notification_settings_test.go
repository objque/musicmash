package api

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/stretchr/testify/assert"
)

func TestAPI_NotificationSettings_Create(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := notifysettings.Create(client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.NoError(t, err)
	settings, err := notifysettings.List(client, testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, settings, 1)
	assert.Equal(t, "telegram", settings[0].Service)
	assert.Equal(t, "chat-id-here", settings[0].Data)
}

func TestAPI_NotificationSettings_Create_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Create(client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.Error(t, err)
}

func TestAPI_NotificationSettings_Update(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "icq",
		Data:     "chat-id-here",
	}))
	assert.NoError(t, db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Update(client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "new-chat-id-here",
	})

	// assert
	assert.NoError(t, err)
	settings, err := notifysettings.List(client, testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, settings, 2)
	assert.Equal(t, "icq", settings[0].Service)
	assert.Equal(t, "chat-id-here", settings[0].Data)
	assert.Equal(t, "telegram", settings[1].Service)
	assert.Equal(t, "new-chat-id-here", settings[1].Data)
}
