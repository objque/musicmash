package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Subscriptions_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureSubscriptionExists(&Subscription{UserID: "objque@me", ArtistName: "skrillex"})

	// assert
	assert.NoError(t, err)
	assert.True(t, DbMgr.IsUserSubscribedForArtist("objque@me", "skrillex"))
}

func TestDB_Subscriptions_FindAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{UserID: "objque@me", ArtistName: "skrillex"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{UserID: "objque@me", ArtistName: "alvin risk"}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{UserID: "jade@dynasty", ArtistName: "rammstein"}))

	// action
	subs, err := DbMgr.FindAllUserSubscriptions("objque@me")

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.Equal(t, "skrillex", subs[0].ArtistName)
	assert.Equal(t, "alvin risk", subs[1].ArtistName)
}
