package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestStores_EnsureExists() {
	// action
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(vars.StoreDeezer))

	// assert
	assert.True(t.T(), DbMgr.IsStoreExists(vars.StoreDeezer))
}

func (t *testDBSuite) TestStores_NotExists() {
	assert.False(t.T(), DbMgr.IsStoreExists(vars.StoreDeezer))
}

func (t *testDBSuite) TestStores_GetAll() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(vars.StoreDeezer))
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	stores, err := DbMgr.GetAllStores()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), stores, 2)
	assert.Equal(t.T(), vars.StoreDeezer, stores[0].Name)
	assert.Equal(t.T(), vars.StoreApple, stores[1].Name)
}
