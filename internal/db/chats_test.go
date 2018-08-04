package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Chat_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureChatExists(&Chat{
		ID: 10000420,
		UserID: "objque@me",
	})

	// assert
	assert.NoError(t, err)
	chatID, err := DbMgr.FindChatByUserID("objque@me")
	assert.NoError(t, err)
	assert.Equal(t, int64(10000420), *chatID)
}
