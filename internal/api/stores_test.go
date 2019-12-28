package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestStores_List() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	stores, err := stores.List(t.client)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), stores, 1)
	assert.Equal(t.T(), vars.StoreApple, stores[0].Name)
}
