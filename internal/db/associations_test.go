package db

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Associations_EnsureArtistExistsInStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreDeezer))

	// action
	err := DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreDeezer, testutil.StoreIDA)

	// assert
	assert.NoError(t, err)
	artists, err := DbMgr.GetAllAssociationsFromStore(testutil.StoreDeezer)
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}

func TestDB_Associations_GetArtistFromStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(testutil.StoreApple))
	assert.NoError(t, DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t, DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDB))

	// action
	artists, err := DbMgr.GetAssociationFromStore(testutil.StoreIDQ, testutil.StoreApple)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
}
