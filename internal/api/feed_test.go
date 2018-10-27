package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Feed_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userName = "objque@me"
	assert.NoError(t, db.DbMgr.EnsureUserExists(userName))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed", server.URL, userName))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Feed_Get_WithQuery(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userName = "objque@me"
	assert.NoError(t, db.DbMgr.EnsureUserExists(userName))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "S.P.Y",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, db.DbMgr.SubscribeUserForArtists(userName, []string{"skrillex", "S.P.Y"}))

	// action
	const layout = "2006-01-02"
	since := time.Now().UTC().Add(-time.Hour * 24 * 2) // two days ago
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed?since=%s", server.URL, userName, since.Format(layout)))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Feed_Get_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(fmt.Sprintf("%s/objque@me/feed", server.URL))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
