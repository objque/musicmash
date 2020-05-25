package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSearch_Do() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects}))

	// action
	result, err := Mgr.Search("skri")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), result.Artists, 1)
	assert.Equal(t.T(), vars.ArtistSkrillex, result.Artists[0].Name)
}

func (t *testDBSuite) TestSearch_Do_Empty() {
	// action
	result, err := Mgr.Search(vars.ArtistArchitects)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), result.Artists, 0)
}
