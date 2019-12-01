package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func (t *testApiSuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDQ}))

	// action
	err := subscriptions.Create(t.client, testutil.UserObjque, []int64{
		testutil.StoreIDQ, testutil.StoreIDW,
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), subs[0].ArtistID)
}

func (t *testApiSuite) TestSubscriptions_List() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDW}))
	assert.NoError(t.T(), db.DbMgr.SubscribeUser(testutil.UserObjque, []int64{
		testutil.StoreIDQ, testutil.StoreIDW,
	}))

	// action
	err := subscriptions.Delete(t.client, testutil.UserObjque, []int64{testutil.StoreIDW})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutil.StoreIDQ), subs[0].ArtistID)
}
