package db

import (
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), testutil.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), subs[0].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetArtists() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserBot, []int64{testutil.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ, testutil.StoreIDW}))

	// action
	subs, err := DbMgr.GetArtistsSubscriptions([]int64{testutil.StoreIDW})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	assert.Equal(t.T(), testutil.UserBot, subs[0].UserName)
	assert.Equal(t.T(), int64(testutil.StoreIDW), subs[0].ArtistID)
	assert.Equal(t.T(), testutil.UserObjque, subs[1].UserName)
	assert.Equal(t.T(), int64(testutil.StoreIDW), subs[1].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutil.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), subs[0])
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUserSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutil.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ}))
	subs, err := DbMgr.GetUserSubscriptions(testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testutil.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), subs[0].ArtistID)

	// action
	err = DbMgr.UnSubscribeUser(testutil.UserObjque, []int64{testutil.StoreIDQ})

	// assert
	assert.NoError(t.T(), err)
	subs, err = DbMgr.GetUserSubscriptions(testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}
