package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestArtists_Create() {
	// action
	artist := &artists.Artist{Name: vars.ArtistArchitects}
	err := artists.Create(t.client, artist)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.ArtistArchitects, artist.Name)
	assert.Equal(t.T(), int64(1), artist.ID)
}

func (t *testAPISuite) TestArtists_Get() {
	// arrange
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistArchitects}))

	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), int64(1), artist.ID)
	assert.Equal(t.T(), vars.ArtistArchitects, artist.Name)
}

func (t *testAPISuite) TestArtists_Get_NotFound() {
	// action
	artist, err := artists.Get(t.client, 1, nil)

	// assert
	assert.Error(t.T(), err)
	assert.Nil(t.T(), artist)
}
