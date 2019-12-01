package db

import (
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestStoreType() {
	// action
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutil.StoreDeezer))

	// assert
	assert.True(t.T(), DbMgr.IsStoreExists(testutil.StoreDeezer))
}

func (t *testDBSuite) TestStoreType_NotExists() {
	assert.False(t.T(), DbMgr.IsStoreExists(testutil.StoreDeezer))
}

func (t *testDBSuite) TestStoreType_GetAll() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutil.StoreDeezer))
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	stores, err := DbMgr.GetAllStores()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), stores, 2)
	assert.Equal(t.T(), testutil.StoreDeezer, stores[0].Name)
	assert.Equal(t.T(), testutil.StoreApple, stores[1].Name)
}
