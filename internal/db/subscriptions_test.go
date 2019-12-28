package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutils.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), testutils.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), subs[0].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetArtists() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserBot, []int64{testutils.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ, testutils.StoreIDW}))

	// action
	subs, err := DbMgr.GetArtistsSubscriptions([]int64{testutils.StoreIDW})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	assert.Equal(t.T(), testutils.UserBot, subs[0].UserName)
	assert.Equal(t.T(), int64(testutils.StoreIDW), subs[0].ArtistID)
	assert.Equal(t.T(), testutils.UserObjque, subs[1].UserName)
	assert.Equal(t.T(), int64(testutils.StoreIDW), subs[1].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(testutils.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutils.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), subs[0])
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUserSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(testutils.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ}))
	subs, err := DbMgr.GetUserSubscriptions(testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testutils.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), subs[0].ArtistID)

	// action
	err = DbMgr.UnSubscribeUser(testutils.UserObjque, []int64{testutils.StoreIDQ})

	// assert
	assert.NoError(t.T(), err)
	subs, err = DbMgr.GetUserSubscriptions(testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}
