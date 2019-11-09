package api

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Subscriptions_Create(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDQ}))

	// action
	err := subscriptions.Create(client, testutil.UserObjque, []int64{
		testutil.StoreIDQ, testutil.StoreIDW,
	})

	// assert
	assert.NoError(t, err)
	subs, err := subscriptions.List(client, testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), subs[0].ArtistID)
}

func TestAPI_Subscriptions_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDQ}))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{ID: testutil.StoreIDW}))
	assert.NoError(t, db.DbMgr.SubscribeUser(testutil.UserObjque, []int64{
		testutil.StoreIDQ, testutil.StoreIDW,
	}))

	// action
	err := subscriptions.Delete(client, testutil.UserObjque, []int64{testutil.StoreIDW})

	// assert
	assert.NoError(t, err)
	subs, err := subscriptions.List(client, testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), subs[0].ArtistID)
}
