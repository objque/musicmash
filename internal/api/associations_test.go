package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/associations"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestAssociations_Add() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	association := &associations.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := associations.Create(t.client, association)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.StoreIDA, association.StoreID)
	assert.Equal(t.T(), vars.StoreApple, association.StoreName)
	assert.Equal(t.T(), int64(1), association.ArtistID)
}

func (t *testAPISuite) TestAssociations_ArtistNotFound() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))

	// action
	association := &associations.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := associations.Create(t.client, association)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestAssociations_AlreadyAssociated() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), db.DbMgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), db.DbMgr.EnsureAssociationExists(1, vars.StoreApple, vars.StoreIDA))

	// action
	association := &associations.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := associations.Create(t.client, association)

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestAssociations_StoreNotFound() {
	// action
	association := &associations.Association{ArtistID: 1, StoreName: vars.StoreApple, StoreID: vars.StoreIDA}
	err := associations.Create(t.client, association)

	// assert
	assert.Error(t.T(), err)
}
