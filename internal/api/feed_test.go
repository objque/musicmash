package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Feed_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userID = "objque@me"
	assert.NoError(t, db.DbMgr.EnsureUserExists(userID))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(&db.Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(&db.Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))

	// action
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed", server.URL, userID))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestAPI_Feed_Get_WithQuery(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userID = "objque@me"
	assert.NoError(t, db.DbMgr.EnsureUserExists(userID))
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, db.DbMgr.EnsureReleaseExists(&db.Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(&db.Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, db.DbMgr.EnsureSubscriptionExists(&db.Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))

	// action
	const layout = "2006-01-02"
	since := time.Now().UTC().Add(-time.Hour * 24 * 2) // two days ago
	resp, err := http.Get(fmt.Sprintf("%s/%s/feed?since=%s", server.URL, userID, since.Format(layout)))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestAPI_Feed_Get_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(fmt.Sprintf("%s/objque@me/feed", server.URL))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}
