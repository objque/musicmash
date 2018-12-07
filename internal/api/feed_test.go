package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Feed_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed", server.URL, testutil.UserObjque))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Feed_Get_WithQuery(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDA,
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: testutil.ArtistSPY,
		StoreName:  testutil.StoreApple,
		StoreID:    testutil.StoreIDB,
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 24 * 2) // two days ago
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed?since=%s", server.URL, testutil.UserObjque, since.Format(testutil.DateYYYYHHMM)))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Feed_Get_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed", server.URL, testutil.UserObjque))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
