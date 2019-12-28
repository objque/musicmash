package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(vars.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(vars.StoreIDQ), subs[0].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetArtists() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserBot, []int64{vars.StoreIDW}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ, vars.StoreIDW}))

	// action
	subs, err := DbMgr.GetArtistsSubscriptions([]int64{vars.StoreIDW})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	assert.Equal(t.T(), vars.UserBot, subs[0].UserName)
	assert.Equal(t.T(), int64(vars.StoreIDW), subs[0].ArtistID)
	assert.Equal(t.T(), vars.UserObjque, subs[1].UserName)
	assert.Equal(t.T(), int64(vars.StoreIDW), subs[1].ArtistID)
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUser() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))

	// action
	subs, err := DbMgr.GetUserSubscriptions(vars.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGetSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(vars.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), subs[0])
}

func (t *testDBSuite) TestSubscriptions_Get_ForAnotherUserSimple() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))

	// action
	subs, err := DbMgr.GetSimpleUserSubscriptions(vars.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))
	subs, err := DbMgr.GetUserSubscriptions(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(vars.StoreIDQ), subs[0].ArtistID)

	// action
	err = DbMgr.UnSubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ})

	// assert
	assert.NoError(t.T(), err)
	subs, err = DbMgr.GetUserSubscriptions(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}
