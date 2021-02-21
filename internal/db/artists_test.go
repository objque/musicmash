package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestArtists_EnsureExists() {
	// action
	err := Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex})

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestArtists_GetAll() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))

	// action
	artists, err := Mgr.GetAllArtists()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestArtists_Get() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 1, Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 2, Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 3, Name: vars.ArtistSPY}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 4, Name: vars.ArtistWildways}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 5, Name: vars.ArtistRitaOra}))

	// action
	artist, err := Mgr.GetArtist(1)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), vars.ArtistSkrillex, artist.Name)
}
