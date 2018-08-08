package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Users_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureUserExists("objque@me")

	// assert
	assert.NoError(t, err)
	user, err := DbMgr.FindUserByID("objque@me")
	assert.NoError(t, err)
	assert.Equal(t, "objque@me", user.ID)
}

func TestDB_Users_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureUserExists("objque@me"))
	assert.NoError(t, DbMgr.EnsureUserExists("jade@abuse"))

	// action
	users, err := DbMgr.GetAllUsers()

	// assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "objque@me", users[0].ID)
	assert.Equal(t, "jade@abuse", users[1].ID)
}
