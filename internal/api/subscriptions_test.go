package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{ID: vars.StoreIDQ}))

	// action
	err := subscriptions.Create(t.client, vars.UserObjque, []*subscriptions.Subscription{
		{ArtistID: vars.StoreIDQ},
		{ArtistID: vars.StoreIDW},
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{ID: vars.StoreIDQ}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{ID: vars.StoreIDW}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{
		vars.StoreIDQ, vars.StoreIDW,
	}))

	// action
	err := subscriptions.Delete(t.client, vars.UserObjque, []*subscriptions.Subscription{
		{ArtistID: vars.StoreIDW},
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, vars.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), subs[0].ArtistID)
}
