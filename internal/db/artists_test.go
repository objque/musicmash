package db

import (
	"github.com/musicmash/musicmash/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestArtist_EnsureExists() {
	// action
	err := DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistSkrillex})

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestArtists_GetAll() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistSkrillex}))

	// action
	artists, err := DbMgr.GetAllArtists()

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), artists, 1)
}

func (t *testDBSuite) TestArtists_Search() {
	// arrange
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistSkrillex, Followers: 100}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistArchitects, Followers: 250}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistSPY}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistWildways, Followers: 50}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{Name: testutils.ArtistRitaOra, Followers: 90}))
	want := []struct {
		SearchText string
		Artists    []string
	}{
		{SearchText: "il", Artists: []string{testutils.ArtistSkrillex, testutils.ArtistWildways}},
		{SearchText: testutils.ArtistSkrillex, Artists: []string{testutils.ArtistSkrillex}},
		{SearchText: "a", Artists: []string{testutils.ArtistArchitects, testutils.ArtistRitaOra, testutils.ArtistWildways}},
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
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: testutils.ArtistSkrillex, Followers: 100}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: testutils.ArtistArchitects, Followers: 250}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: testutils.ArtistSPY}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: testutils.ArtistWildways, Followers: 50}))
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: testutils.ArtistRitaOra, Followers: 90}))

	// action
	artist, err := DbMgr.GetArtist(1)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), testutils.ArtistSkrillex, artist.Name)
	assert.Equal(t.T(), uint(100), artist.Followers)
}
