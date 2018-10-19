package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Subscriptions_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureSubscriptionExists("objque@me", "skrillex")

	// assert
	assert.NoError(t, err)
	assert.True(t, DbMgr.IsUserSubscribedForArtist("objque@me", "skrillex"))
}

func TestDB_Subscriptions_FindAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "skrillex"))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "alvin risk"))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("jade@dynasty", "rammstein"))

	// action
	subs, err := DbMgr.FindAllUserSubscriptions("objque@me")

	// assert
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
}

func TestDB_Subscriptions_SubscribeUserForArtists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	artists := []string{"Skrillex", "David Guetta", "Deftones", "Depeche Mode"}
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "Skrillex"))

	// action
	err := DbMgr.SubscribeUserForArtists("objque@me", artists)

	// assert
	assert.NoError(t, err)
	subs, err := DbMgr.FindAllUserSubscriptions("objque@me")
	assert.NoError(t, err)
	assert.Len(t, subs, 4)
}

func TestDB_Subscriptions_UnsubscribeUserFromArtists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "Skrillex"))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "Calvin Risk"))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("objque@me", "AC/DC"))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists("mike@wels", "AC/DC"))

	// action
	err := DbMgr.UnsubscribeUserFromArtists("objque@me", []string{"Calvin Risk"})

	// assert
	assert.NoError(t, err)

	subs, err := DbMgr.FindAllUserSubscriptions("mike@wels")
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, "AC/DC", subs[0].ArtistName)

	subs, err = DbMgr.FindAllUserSubscriptions("objque@me")
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.Equal(t, "AC/DC", subs[0].ArtistName)
	assert.Equal(t, "Skrillex", subs[1].ArtistName)
}
