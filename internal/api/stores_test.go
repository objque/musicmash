package api

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Stores_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureStoreExists(testutil.StoreApple))

	// action
	stores, err := stores.List(client)

	// assert
	assert.NoError(t, err)
	assert.Len(t, stores, 1)
	assert.Equal(t, testutil.StoreApple, stores[0].Name)
}
