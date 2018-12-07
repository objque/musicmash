package db

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Users_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureUserExists(testutil.UserObjque)

	// assert
	assert.NoError(t, err)
	user, err := DbMgr.FindUserByName(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Equal(t, testutil.UserObjque, user.Name)
}

func TestDB_Users_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserBot))

	// action
	users, err := DbMgr.GetAllUsers()

	// assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, testutil.UserObjque, users[0].Name)
	assert.Equal(t, testutil.UserBot, users[1].Name)
}

func TestDB_Users_GetUsersWithReleases(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserBot))
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserTest))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserTest, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{ArtistName: testutil.ArtistArchitects, CreatedAt: now}))

	// action
	users, err := DbMgr.GetUsersWithReleases(now)

	// assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.EqualValues(t, []string{testutil.UserObjque, testutil.UserTest}, users)
}

func TestDB_Users_GetUsersWithReleases_NoReleases_ForProvidedHour(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now().UTC()
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserObjque))
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserBot))
	assert.NoError(t, DbMgr.EnsureUserExists(testutil.UserTest))
	assert.NoError(t, DbMgr.EnsureArtistExists(testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserObjque, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(testutil.UserTest, testutil.ArtistArchitects))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{ArtistName: testutil.ArtistArchitects, CreatedAt: now.Add(-time.Hour)}))

	// action
	users, err := DbMgr.GetUsersWithReleases(now)

	// assert
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
