package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Subscriptions_SubscribeAndGetUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserObjque)

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, testutil.UserObjque, subs[0].UserName)
	assert.Equal(t, int64(testutil.StoreIDQ), subs[0].ArtistID)
}

func TestDB_Subscriptions_SubscribeAndGetArtists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDW}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserBot, []int64{testutil.StoreIDW}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ, testutil.StoreIDW}))

	// action
	subs, err := DbMgr.GetArtistsSubscriptions([]int64{testutil.StoreIDW})

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.Equal(t, testutil.UserBot, subs[0].UserName)
	assert.Equal(t, int64(testutil.StoreIDW), subs[0].ArtistID)
	assert.Equal(t, testutil.UserObjque, subs[1].UserName)
	assert.Equal(t, int64(testutil.StoreIDW), subs[1].ArtistID)
}

func TestDB_Subscriptions_Get_ForAnotherUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserBot)

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 0)
}

func TestDB_Subscriptions_SubscribeAndGetSimple(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutil.UserObjque)

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), subs[0])
}

func TestDB_Subscriptions_Get_ForAnotherUserSimple(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutil.UserBot)

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 0)
}

func TestDB_Subscriptions_UnSubscribe(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Equal(t, testutil.UserObjque, subs[0].UserName)
	assert.Equal(t, int64(testutil.StoreIDQ), subs[0].ArtistID)

	// action
	err = DbMgr.UnSubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ})

	// assert
	assert.NoError(t, err)
	subs, err = DbMgr.GetUserSubscriptions(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 0)
}
