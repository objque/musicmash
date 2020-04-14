package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestAssociations_EnsureArtistExistsInStore() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreDeezer))

	// action
	err := Mgr.EnsureAssociationExists(vars.StoreIDQ, vars.StoreDeezer, vars.StoreIDA)

	// assert
	assert.NoError(t.T(), err)
	artists, err := Mgr.GetAllAssociationsFromStore(vars.StoreDeezer)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestAssociations_GetArtistFromStore() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), Mgr.EnsureAssociationExists(vars.StoreIDQ, vars.StoreApple, vars.StoreIDA))
	assert.NoError(t.T(), Mgr.EnsureAssociationExists(vars.StoreIDQ, vars.StoreApple, vars.StoreIDB))

	// action
	artists, err := Mgr.GetAssociationFromStore(vars.StoreIDQ, vars.StoreApple)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 2)
}
