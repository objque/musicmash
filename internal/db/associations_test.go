package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestAssociations_EnsureArtistExistsInStore() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutils.StoreDeezer))

	// action
	err := DbMgr.EnsureAssociationExists(testutils.StoreIDQ, testutils.StoreDeezer, testutils.StoreIDA)

	// assert
	assert.NoError(t.T(), err)
	artists, err := DbMgr.GetAllAssociationsFromStore(testutils.StoreDeezer)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestAssociations_GetArtistFromStore() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureStoreExists(testutils.StoreApple))
	assert.NoError(t.T(), DbMgr.EnsureAssociationExists(testutils.StoreIDQ, testutils.StoreApple, testutils.StoreIDA))
	assert.NoError(t.T(), DbMgr.EnsureAssociationExists(testutils.StoreIDQ, testutils.StoreApple, testutils.StoreIDB))

	// action
	artists, err := DbMgr.GetAssociationFromStore(testutils.StoreIDQ, testutils.StoreApple)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 2)
}
