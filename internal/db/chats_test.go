package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Chat_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureChatExists(&Chat{ID: 10000420, UserName: "objque@me"})

	// assert
	assert.NoError(t, err)
	chatID, err := DbMgr.FindChatByUserName("objque@me")
	assert.NoError(t, err)
	assert.Equal(t, int64(10000420), *chatID)
}

func TestDB_Chat_GetAllChatsThatSubscribedFor(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// chats
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 1001, UserName: "objque@me"}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 3004, UserName: "another@user"}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 5005, UserName: "friction@user"}))
	// subs
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Skrillex", UserName: "objque@me"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Rammstein", UserName: "another@user"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Rammstein", UserName: "objque@me"}))
	want := []struct {
		Artist     string
		ChatsCount int
	}{
		{ChatsCount: 0, Artist: "Qeen"},
		{ChatsCount: 1, Artist: "Skrillex"},
		{ChatsCount: 2, Artist: "Rammstein"},
	}

	for _, w := range want {
		// action
		chats, err := DbMgr.GetAllChatsThatSubscribedFor(w.Artist)

		// assert
		assert.NoError(t, err)
		assert.Len(t, chats, w.ChatsCount)
	}
}
