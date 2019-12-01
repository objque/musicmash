package db

import (
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestAssociations_EnsureArtistExistsInStore() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutil.StoreDeezer))

	// action
	err := DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreDeezer, testutil.StoreIDA)

	// assert
	assert.NoError(t.T(), err)
	artists, err := DbMgr.GetAllAssociationsFromStore(testutil.StoreDeezer)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestAssociations_GetArtistFromStore() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutil.StoreApple))
	assert.NoError(t.T(), DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDA))
	assert.NoError(t.T(), DbMgr.EnsureAssociationExists(testutil.StoreIDQ, testutil.StoreApple, testutil.StoreIDB))

	// action
	artists, err := DbMgr.GetAssociationFromStore(testutil.StoreIDQ, testutil.StoreApple)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 2)
}
