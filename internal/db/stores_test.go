package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestStores_EnsureExists() {
	// action
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutils.StoreDeezer))

	// assert
	assert.True(t.T(), DbMgr.IsStoreExists(testutils.StoreDeezer))
}

func (t *testDBSuite) TestStores_NotExists() {
	assert.False(t.T(), DbMgr.IsStoreExists(testutils.StoreDeezer))
}

func (t *testDBSuite) TestStores_GetAll() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutils.StoreDeezer))
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutils.StoreApple))

	// action
	stores, err := DbMgr.GetAllStores()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), stores, 2)
	assert.Equal(t.T(), testutils.StoreDeezer, stores[0].Name)
	assert.Equal(t.T(), testutils.StoreApple, stores[1].Name)
}
