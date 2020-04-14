package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestReleasesBatch_InsertMany_1() {
	// arrange
	now := time.Date(2020, 10, 01, 0, 0, 0, 0, time.UTC)
	rels := []*Release{
		{
			CreatedAt: now,
			ArtistID:  vars.StoreIDW,
			Title:     vars.ReleaseArchitectsHollyHell,
			Poster:    vars.PosterMiddle,
			Released:  now.AddDate(0, -1, 0),
			StoreName: vars.StoreApple,
			StoreID:   vars.StoreIDA,
			Type:      vars.ReleaseTypeVideo,
			Explicit:  true,
		}, {
			CreatedAt: now,
			ArtistID:  vars.StoreIDQ,
			Title:     vars.ReleaseArchitectsHollyHell,
			Poster:    vars.PosterMiddle,
			Released:  now.AddDate(0, -5, 0),
			StoreName: vars.StoreSpotify,
			StoreID:   vars.StoreIDB,
			Type:      vars.ReleaseTypeSong,
			Explicit:  false,
		},
	}

	// action
	err := Mgr.InsertBatchNewReleases(rels)

	// assert
	assert.NoError(t.T(), err)
	releases, err := Mgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 2)

	assert.Equal(t.T(), now.Year(), releases[0].CreatedAt.Year())
	assert.Equal(t.T(), now.Day(), releases[0].CreatedAt.Day())
	assert.Equal(t.T(), int64(vars.StoreIDW), releases[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, releases[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, releases[0].Poster)
	assert.Equal(t.T(), now.Year(), releases[0].Released.Year())
	assert.Equal(t.T(), now.Day(), releases[0].Released.Day())
	assert.Equal(t.T(), vars.StoreApple, releases[0].StoreName)
	assert.Equal(t.T(), vars.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, releases[0].Type)
	assert.Equal(t.T(), true, releases[0].Explicit)

	assert.Equal(t.T(), now.Year(), releases[1].CreatedAt.Year())
	assert.Equal(t.T(), now.Day(), releases[1].CreatedAt.Day())
	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, releases[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, releases[1].Poster)
	assert.Equal(t.T(), now.Year(), releases[1].Released.Year())
	assert.Equal(t.T(), now.Day(), releases[1].Released.Day())
	assert.Equal(t.T(), vars.StoreSpotify, releases[1].StoreName)
	assert.Equal(t.T(), vars.StoreIDB, releases[1].StoreID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, releases[1].Type)
	assert.Equal(t.T(), false, releases[1].Explicit)
}

func (t *testDBSuite) TestReleasesBatch_InsertMany_Ignore() {
	// arrange
	now := time.Date(2020, 10, 01, 0, 0, 0, 0, time.UTC)
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(&Release{
		CreatedAt: now,
		ArtistID:  vars.StoreIDW,
		Title:     vars.ReleaseArchitectsHollyHell,
		Poster:    vars.PosterMiddle,
		Released:  now.AddDate(0, -1, 0),
		StoreName: vars.StoreApple,
		StoreID:   vars.StoreIDA,
		Type:      vars.ReleaseTypeVideo,
		Explicit:  true,
	}))
	rels := []*Release{
		{
			CreatedAt: now,
			ArtistID:  vars.StoreIDW,
			Title:     vars.ReleaseArchitectsHollyHell,
			Poster:    vars.PosterMiddle,
			Released:  now.AddDate(0, -1, 0),
			StoreName: vars.StoreApple,
			StoreID:   vars.StoreIDA,
			Type:      vars.ReleaseTypeVideo,
			Explicit:  true,
		}, {
			CreatedAt: now,
			ArtistID:  vars.StoreIDQ,
			Title:     vars.ReleaseArchitectsHollyHell,
			Poster:    vars.PosterMiddle,
			Released:  now.AddDate(0, -5, 0),
			StoreName: vars.StoreSpotify,
			StoreID:   vars.StoreIDB,
			Type:      vars.ReleaseTypeSong,
			Explicit:  false,
		},
	}

	// action
	err := Mgr.InsertBatchNewReleases(rels)

	// assert
	assert.NoError(t.T(), err)
	releases, err := Mgr.GetAllReleases()
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 2)

	assert.Equal(t.T(), now.Year(), releases[0].CreatedAt.Year())
	assert.Equal(t.T(), now.Day(), releases[0].CreatedAt.Day())
	assert.Equal(t.T(), int64(vars.StoreIDW), releases[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, releases[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, releases[0].Poster)
	assert.Equal(t.T(), now.Year(), releases[0].Released.Year())
	assert.Equal(t.T(), now.Day(), releases[0].Released.Day())
	assert.Equal(t.T(), vars.StoreApple, releases[0].StoreName)
	assert.Equal(t.T(), vars.StoreIDA, releases[0].StoreID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, releases[0].Type)
	assert.Equal(t.T(), true, releases[0].Explicit)

	assert.Equal(t.T(), now.Year(), releases[1].CreatedAt.Year())
	assert.Equal(t.T(), now.Day(), releases[1].CreatedAt.Day())
	assert.Equal(t.T(), int64(vars.StoreIDQ), releases[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, releases[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, releases[1].Poster)
	assert.Equal(t.T(), now.Year(), releases[1].Released.Year())
	assert.Equal(t.T(), now.Day(), releases[1].Released.Day())
	assert.Equal(t.T(), vars.StoreSpotify, releases[1].StoreName)
	assert.Equal(t.T(), vars.StoreIDB, releases[1].StoreID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, releases[1].Type)
	assert.Equal(t.T(), false, releases[1].Explicit)
}
