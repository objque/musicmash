package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Chat_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureChatExists(&Chat{ID: 10000420, UserName: testutil.UserObjque})

	// assert
	assert.NoError(t, err)
	chatID, err := DbMgr.FindChatByUserName(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Equal(t, int64(10000420), *chatID)
}

func TestDB_Chat_GetAllChatsThatSubscribedFor(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 10000420, UserName:testutil.UserObjque}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserBot, []int64{testutil.StoreIDW}))

	// action
	chats, err := DbMgr.GetAllChatsThatSubscribedFor(testutil.StoreIDQ)

	// assert
	assert.NoError(t, err)
	assert.Len(t, chats, 1)
	assert.Equal(t, int64(10000420), chats[0].ID)
}
