package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_StoreType(t *testing.T) {
	setup()
	defer teardown()

	// action
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreDeezer))

	// assert
	assert.True(t, DbMgr.IsStoreExists(testutil.StoreDeezer))
}

func TestDB_StoreType_NotExists(t *testing.T) {
	setup()
	defer teardown()

	assert.False(t, DbMgr.IsStoreExists(testutil.StoreDeezer))
}

func TestDB_StoreType_GetAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreDeezer))
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	stores, err := DbMgr.GetAllStores()

	// assert
	assert.NoError(t, err)
	assert.Len(t, stores, 2)
	assert.Equal(t, testutil.StoreDeezer, stores[0].Name)
	assert.Equal(t, testutil.StoreApple, stores[1].Name)
}
