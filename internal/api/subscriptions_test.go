package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutils.StoreIDQ}))

	// action
	err := subscriptions.Create(t.client, testutils.UserObjque, []int64{
		testutils.StoreIDQ, testutils.StoreIDW,
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), subs[0].ArtistID)
}

func (t *testAPISuite) TestSubscriptions_List() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutils.StoreIDQ}))
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutils.StoreIDW}))
	assert.NoError(t.T(), db.DbMgr.SubscribeUser(testutils.UserObjque, []int64{
		testutils.StoreIDQ, testutils.StoreIDW,
	}))

	// action
	err := subscriptions.Delete(t.client, testutils.UserObjque, []int64{testutils.StoreIDW})

	// assert
	assert.NoError(t.T(), err)
	subs, err := subscriptions.List(t.client, testutils.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	assert.Equal(t.T(), int64(testutils.StoreIDQ), subs[0].ArtistID)
}
