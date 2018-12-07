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
