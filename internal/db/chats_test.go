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
	// chats
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 1001, UserName: testutil.UserObjque}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 3004, UserName: testutil.UserBot}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 5005, UserName: testutil.UserTest}))
	// subs
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserBot, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistArchitects))
	want := []struct {
		Artist     string
		ChatsCount int
	}{
		{ChatsCount: 0, Artist: testutil.ArtistSPY},
		{ChatsCount: 1, Artist: testutil.ArtistSkrillex},
		{ChatsCount: 2, Artist: testutil.ArtistArchitects},
	}

	for _, w := range want {
		// action
		chats, err := DbMgr.GetAllChatsThatSubscribedFor(w.Artist)

		// assert
		assert.NoError(t, err)
		assert.Len(t, chats, w.ChatsCount)
	}
}
