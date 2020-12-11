package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleases_EnsureExists() {
	// action
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreDeezer))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	err := Mgr.EnsureReleaseExists(&Release{
		ArtistID:  1,
		SpotifyID: vars.StoreApple,
		Explicit:  true,
	})

	// assert
	assert.NoError(t.T(), err)
	releases, err := Mgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.True(t.T(), releases[0].Explicit)
}

func (t *testDBSuite) TestReleases_FindReleases() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSPY}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  1,
		SpotifyID: vars.StoreIDA,
		Title:     vars.ArtistAlgorithm,
		Explicit:  true,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		ArtistID:  2,
		SpotifyID: vars.StoreIDB,
	}))

	// action
	releases, err := Mgr.FindReleases(1, vars.ArtistAlgorithm)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	assert.Equal(t.T(), int64(1), releases[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistAlgorithm, releases[0].Title)
	assert.True(t.T(), releases[0].Explicit)
}
