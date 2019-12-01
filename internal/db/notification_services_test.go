package db

import (
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestNotificationService_Exists() {
	// action
	assert.NoError(t.T(), DbMgr.EnsureNotificationServiceExists("email"))

	// assert
	assert.True(t.T(), DbMgr.IsNotificationServiceExists("email"))
}

func (t *testDBSuite) TestNotificationService_NotExists() {
	assert.False(t.T(), DbMgr.IsNotificationServiceExists("email"))
}
