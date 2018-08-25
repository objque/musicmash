package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Store(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const store = "spotify"

	// action
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore(store, "654321", 001))

	// assert
	assert.True(t, DbMgr.IsReleaseExistsInStore(store, "654321"))
}

func TestDB_Store_NotExists(t *testing.T) {
	setup()
	defer teardown()

	assert.False(t, DbMgr.IsReleaseExistsInStore("spotify", "654321"))
}
