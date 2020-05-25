package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/search"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestSearch_Do() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))

	// action
	result, err := search.Do(t.client, "skri")

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), result.Artists, 1)
	assert.Equal(t.T(), vars.ArtistSkrillex, result.Artists[0].Name)
}

func (t *testAPISuite) TestSearch_Do_Empty() {
	// action
	result, err := search.Do(t.client, vars.ArtistArchitects)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), result.Artists, 0)
}
