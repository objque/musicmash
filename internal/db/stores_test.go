package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_StoreType(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const store = "spotify"

	// action
	assert.NoError(t, DbMgr.EnsureStoreExists(store))

	// assert
	assert.True(t, DbMgr.IsStoreExists(store))
}

func TestDB_StoreType_NotExists(t *testing.T) {
	setup()
	defer teardown()

	assert.False(t, DbMgr.IsStoreExists("deezer"))
}
