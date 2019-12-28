package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestArtist_EnsureExists() {
	// action
	err := DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex})

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestArtists_GetAll() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))

	// action
	artists, err := DbMgr.GetAllArtists()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestArtists_Search() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex, Followers: 100}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects, Followers: 250}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistWildways, Followers: 50}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistRitaOra, Followers: 90}))
	want := []struct {
		SearchText string
		Artists    []string
	}{
		{SearchText: "il", Artists: []string{vars.ArtistSkrillex, vars.ArtistWildways}},
		{SearchText: vars.ArtistSkrillex, Artists: []string{vars.ArtistSkrillex}},
		{SearchText: "a", Artists: []string{vars.ArtistArchitects, vars.ArtistRitaOra, vars.ArtistWildways}},
	}

	for i := range want {
		// action
		artists, err := DbMgr.SearchArtists(want[i].SearchText)

		// assert
		assert.NoError(t.T(), err)
		assert.Len(t.T(), artists, len(want[i].Artists))
		for i, wantName := range want[i].Artists {
			assert.Equal(t.T(), wantName, artists[i].Name)
		}
	}
}

func (t *testDBSuite) TestArtists_GetWithFullInfo() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: vars.ArtistSkrillex, Followers: 100}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: vars.ArtistArchitects, Followers: 250}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: vars.ArtistSPY}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: vars.ArtistWildways, Followers: 50}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: vars.ArtistRitaOra, Followers: 90}))

	// action
	artist, err := DbMgr.GetArtist(1)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), vars.ArtistSkrillex, artist.Name)
	assert.Equal(t.T(), uint(100), artist.Followers)
}
