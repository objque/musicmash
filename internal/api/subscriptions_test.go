package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAPI_SubscribeUser_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Post(fmt.Sprintf("%s/%s/subscriptions", server.URL, testutil.UserObjque), "application/json", nil)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_UnsubscribeUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistSkrillex))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistWildways))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistArchitects))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(testutil.UserBot, testutil.ArtistArchitects))

	// action
	artists := []string{testutil.ArtistWildways}
	buffer, _ := json.Marshal(&artists)
	resp, err := httpDelete(fmt.Sprintf("%s/%s/subscriptions", server.URL, testutil.UserObjque), bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	subs, err := db.DbMgr.FindAllUserSubscriptions(testutil.UserBot)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, testutil.ArtistArchitects, subs[0].ArtistName)

	subs, err = db.DbMgr.FindAllUserSubscriptions(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.Equal(t, testutil.ArtistArchitects, subs[0].ArtistName)
	assert.Equal(t, testutil.ArtistSkrillex, subs[1].ArtistName)
}

func TestAPI_UnsubscribeUser_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := httpDelete(fmt.Sprintf("%s/%s/subscriptions", server.URL, testutil.UserObjque), nil)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_Subscriptions_Get_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/subscriptions", server.URL, testutil.UserObjque))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_Subscriptions_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserBot))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(testutil.UserBot, []string{testutil.ArtistArchitects}))

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/subscriptions", server.URL, testutil.UserObjque))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	subs := []*db.Subscription{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&subs))
	assert.Len(t, subs, 2)
	assert.Equal(t, testutil.ArtistSPY, subs[0].ArtistName)
	assert.Equal(t, testutil.ArtistSkrillex, subs[1].ArtistName)
}
