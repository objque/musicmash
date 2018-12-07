package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Subscriptions_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex)

	// assert
	assert.NoError(t, err)
	assert.True(t, DbMgr.IsUserSubscribedForArtist(testutil.UserObjque, testutil.ArtistSkrillex))
}

func TestDB_Subscriptions_FindAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserBot, testutil.ArtistWildways))

	// action
	subs, err := DbMgr.FindAllUserSubscriptions(testutil.UserObjque)

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
}

func TestDB_Subscriptions_SubscribeUserForArtists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	artists := []string{testutil.ArtistSkrillex, testutil.ArtistArchitects, testutil.ArtistWildways, testutil.ArtistSPY}
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex))

	// action
	err := DbMgr.SubscribeUserForArtists(testutil.UserObjque, artists)

	// assert
	assert.NoError(t, err)
	subs, err := DbMgr.FindAllUserSubscriptions(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 4)
}

func TestDB_Subscriptions_UnsubscribeUserFromArtists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistWildways))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSPY))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserBot, testutil.ArtistSPY))

	// action
	err := DbMgr.UnsubscribeUserFromArtists(testutil.UserObjque, []string{testutil.ArtistWildways})

	// assert
	assert.NoError(t, err)

	subs, err := DbMgr.FindAllUserSubscriptions(testutil.UserBot)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, testutil.ArtistSPY, subs[0].ArtistName)

	subs, err = DbMgr.FindAllUserSubscriptions(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.Equal(t, testutil.ArtistSPY, subs[0].ArtistName)
	assert.Equal(t, testutil.ArtistSkrillex, subs[1].ArtistName)
}
