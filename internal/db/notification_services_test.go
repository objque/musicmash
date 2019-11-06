package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_NotificationService_Exists(t *testing.T) {
	setup()
	defer teardown()

	// action
	assert.NoError(t, DbMgr.EnsureNotificationServiceExists("email"))

	// assert
	assert.True(t, DbMgr.IsNotificationServiceExists("email"))
}

func TestDB_NotificationService_NotExists(t *testing.T) {
	setup()
	defer teardown()

	assert.False(t, DbMgr.IsNotificationServiceExists("email"))
}
