package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/stretchr/testify/assert"
)

func (t *testApiSuite) TestStores_List() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	stores, err := stores.List(t.client)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), stores, 1)
	assert.Equal(t.T(), testutil.StoreApple, stores[0].Name)
}
