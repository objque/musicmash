package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestNotificationSettings_Create() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureNotificationServiceExists("telegram"))

	// action
	err := notifysettings.Create(t.client, vars.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := notifysettings.List(t.client, vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 1)
	assert.Equal(t.T(), "telegram", settings[0].Service)
	assert.Equal(t.T(), "chat-id-here", settings[0].Data)
}

func (t *testAPISuite) TestNotificationSettings_Create_AlreadyExists() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureNotificationServiceExists("telegram"))
	assert.NoError(t.T(), db.Mgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: vars.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Create(t.client, vars.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestNotificationSettings_Update() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureNotificationServiceExists("telegram"))
	assert.NoError(t.T(), db.Mgr.EnsureNotificationServiceExists("icq"))
	assert.NoError(t.T(), db.Mgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: vars.UserObjque,
		Service:  "icq",
		Data:     "chat-id-here",
	}))
	assert.NoError(t.T(), db.Mgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: vars.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Update(t.client, vars.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "new-chat-id-here",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := notifysettings.List(t.client, vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 2)
	assert.Equal(t.T(), "icq", settings[0].Service)
	assert.Equal(t.T(), "chat-id-here", settings[0].Data)
	assert.Equal(t.T(), "telegram", settings[1].Service)
	assert.Equal(t.T(), "new-chat-id-here", settings[1].Data)
}
