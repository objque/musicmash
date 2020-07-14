package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))

	// action
	err := Mgr.CreateSubscription(&Subscription{
		UserName: vars.UserObjque,
		ArtistID: 1,
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGet() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))

	// action
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_SubscribeAndGet_Empty() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))

	// action
	subs, err := Mgr.GetUserSubscriptions(vars.UserBot)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)

	// action
	err = Mgr.UnSubscribeUser(vars.UserObjque, []int64{1})

	// assert
	assert.NoError(t.T(), err)
	subs, err = Mgr.GetUserSubscriptions(vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}
