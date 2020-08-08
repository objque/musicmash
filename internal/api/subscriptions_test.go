package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))

	// action
	err := subscriptions.Create(t.client, vars.UserObjque, []*subscriptions.Subscription{
		{ArtistID: 1},
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, vars.UserObjque, nil)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2,
	}))

	// action
	subs, err := subscriptions.List(t.client, vars.UserObjque, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), int64(2), subs[1].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List_Limit() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2,
	}))

	// action
	var limit uint64 = 1
	subs, err := subscriptions.List(t.client, vars.UserObjque, &subscriptions.GetOptions{Limit: &limit})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List_Offset() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2,
	}))

	// action
	var offset uint64 = 1
	subs, err := subscriptions.List(t.client, vars.UserObjque, &subscriptions.GetOptions{Offset: &offset})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(2), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List_Limit_Offset() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2, 3,
	}))

	// action
	var (
		limit  uint64 = 1
		offset uint64 = 1
	)
	subs, err := subscriptions.List(t.client, vars.UserObjque, &subscriptions.GetOptions{
		Limit:  &limit,
		Offset: &offset,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(2), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List_SortType() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2,
	}))

	// action
	const sortType = "desc"
	subs, err := subscriptions.List(t.client, vars.UserObjque, &subscriptions.GetOptions{SortType: sortType})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	assert.Equal(t.T(), int64(2), subs[0].ArtistID)
	assert.Equal(t.T(), int64(1), subs[1].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		1, 2,
	}))

	// action
	err := subscriptions.Delete(t.client, vars.UserObjque, []*subscriptions.Subscription{
		{ArtistID: 2},
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, vars.UserObjque, nil)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
}
