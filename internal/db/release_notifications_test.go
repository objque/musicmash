package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleaseNotifications() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	// create artist
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	// subscribe users
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserBot, []int64{1}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2}))
	// fill releases
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC().AddDate(0, -1, 0),
		ArtistID:  1,
		Poster:    vars.PosterSimple,
		Title:     vars.ReleaseSkrillexHumble,
		Released:  time.Now().UTC().AddDate(0, -1, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output",
		Type:      vars.ReleaseTypeAlbum,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC().AddDate(0, 0, -15),
		ArtistID:  1,
		Poster:    vars.PosterGiant,
		Title:     vars.ReleaseRitaOraLouder,
		Released:  time.Now().UTC().AddDate(0, 0, -15),
		StoreName: vars.StoreApple,
		StoreID:   "this-oldest-release-wont-be-in-output-as-previous",
		Explicit:  true,
		Type:      vars.ReleaseTypeAlbum,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC().AddDate(0, 10, 0),
		ArtistID:  1,
		Poster:    vars.PosterMiddle,
		Title:     vars.ReleaseArchitectsNaySayer,
		Released:  time.Now().UTC().AddDate(0, 10, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-future-release-have-to-be-in-output",
		Type:      vars.ReleaseTypeAlbum,
	}))
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: time.Now().UTC().AddDate(1, 0, 0),
		ArtistID:  2,
		Poster:    vars.PosterSmall,
		Title:     vars.ReleaseSkrillexHumble,
		Released:  time.Now().UTC().AddDate(1, 0, 0),
		StoreName: vars.StoreApple,
		StoreID:   "this-future-release-have-to-be-in-output-as-previous",
		Type:      vars.ReleaseTypeAlbum,
		Explicit:  true,
	}))

	// action
	releases, err := Mgr.GetReleaseNotifications(time.Now().UTC())

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 3)

	release := releases[0]
	assert.Equal(t.T(), vars.UserBot, release.UserName)
	assert.Equal(t.T(), int64(1), release.ArtistID)
	assert.Equal(t.T(), vars.ArtistArchitects, release.ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, release.Poster)
	assert.Equal(t.T(), "this-future-release-have-to-be-in-output", *release.ItunesID)
	assert.Empty(t.T(), release.SpotifyID)
	assert.Empty(t.T(), release.DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, release.Type)
	assert.False(t.T(), release.Explicit)

	release = releases[1]
	assert.Equal(t.T(), vars.UserObjque, release.UserName)
	assert.Equal(t.T(), int64(1), release.ArtistID)
	assert.Equal(t.T(), vars.ArtistArchitects, release.ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, release.Poster)
	assert.Equal(t.T(), "this-future-release-have-to-be-in-output", *release.ItunesID)
	assert.Empty(t.T(), release.SpotifyID)
	assert.Empty(t.T(), release.DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, release.Type)
	assert.False(t.T(), release.Explicit)

	release = releases[2]
	assert.Equal(t.T(), vars.UserObjque, release.UserName)
	assert.Equal(t.T(), int64(2), release.ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, release.ArtistName)
	assert.Equal(t.T(), vars.PosterSmall, release.Poster)
	assert.Equal(t.T(), "this-future-release-have-to-be-in-output-as-previous", *release.ItunesID)
	assert.Empty(t.T(), release.SpotifyID)
	assert.Empty(t.T(), release.DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, release.Type)
	assert.True(t.T(), release.Explicit)
}
