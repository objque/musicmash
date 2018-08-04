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
		ID:     10000420,
		UserID: "objque@me",
	})

	// assert
	assert.NoError(t, err)
	chatID, err := DbMgr.FindChatByUserID("objque@me")
	assert.NoError(t, err)
	assert.Equal(t, int64(10000420), *chatID)
}

func TestDB_Chat_GetAllChatsThatSubscribedFor(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// chats
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 1001, UserID: "objque@me"}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 3004, UserID: "another@user"}))
	assert.NoError(t, DbMgr.EnsureChatExists(&Chat{ID: 5005, UserID: "friction@user"}))
	// subs
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Skrillex", UserID: "objque@me"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Rammstein", UserID: "another@user"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{ArtistName: "Rammstein", UserID: "objque@me"}))
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
