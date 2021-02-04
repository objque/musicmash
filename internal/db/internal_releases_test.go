package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/internal/utils/ptr"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) fillRelease(release *Release) {
	assert.NoError(t.T(), Mgr.EnsureReleaseExists(release))
}

func (t *testDBSuite) prepareReleasesTestCase() {
	r := time.Now().UTC()
	monthAgo := r.AddDate(0, -1, 0)
	yearAgo := r.AddDate(-1, 0, 0)
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreDeezer))
	assert.NoError(t.T(), Mgr.EnsureStoreExists(vars.StoreSpotify))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistAlgorithm, ID: 1}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex, ID: 2}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))
	// first release
	t.fillRelease(&Release{ArtistID: 1, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, SpotifyID: "3000", Type: vars.ReleaseTypeAlbum, Explicit: true, TracksCount: 10, DurationMs: 25})
	// second release
	t.fillRelease(&Release{ArtistID: 1, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterMiddle, Released: monthAgo, SpotifyID: "1100", Type: vars.ReleaseTypeSong, Explicit: false, TracksCount: 10, DurationMs: 25})
	// third release from another artist
	t.fillRelease(&Release{ArtistID: 2, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: yearAgo, SpotifyID: "1110", Type: vars.ReleaseTypeVideo, Explicit: true, TracksCount: 10, DurationMs: 25})
}

func (t *testDBSuite) TestInternalReleases_Get_All() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{SortType: "desc"})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 3)
	// releases are sort by release date desc
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)

	assert.NotEmpty(t.T(), rels[1].ID)
	assert.Equal(t.T(), int64(1), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", rels[1].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.Equal(t.T(), int32(10), rels[1].TracksCount)
	assert.Equal(t.T(), int64(25), rels[1].DurationMs)
	assert.False(t.T(), rels[1].IsExplicit)

	assert.NotEmpty(t.T(), rels[2].ID)
	assert.Equal(t.T(), int64(2), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[2].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[2].Poster)
	assert.Equal(t.T(), "1110", rels[2].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[2].Type)
	assert.Equal(t.T(), int32(10), rels[2].TracksCount)
	assert.Equal(t.T(), int64(25), rels[2].DurationMs)
	assert.True(t.T(), rels[2].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_All_ChangeSortType() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		SortType: releases.SortTypeASC,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 3)
	// releases are sort by release date ASC!
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)

	assert.NotEmpty(t.T(), rels[1].ID)
	assert.Equal(t.T(), int64(1), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", rels[1].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.Equal(t.T(), int32(10), rels[1].TracksCount)
	assert.Equal(t.T(), int64(25), rels[1].DurationMs)
	assert.False(t.T(), rels[1].IsExplicit)

	assert.NotEmpty(t.T(), rels[2].ID)
	assert.Equal(t.T(), int64(1), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[2].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[2].Poster)
	assert.Equal(t.T(), "3000", rels[2].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[2].Type)
	assert.Equal(t.T(), int32(10), rels[2].TracksCount)
	assert.Equal(t.T(), int64(25), rels[2].DurationMs)
	assert.True(t.T(), rels[2].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByLimit() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Limit:    ptr.Uint(1),
		SortType: "desc",
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByArtistID() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		ArtistID: ptr.Int(2),
		SortType: "desc",
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByArtistID_ReleaseType() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		ArtistID:    ptr.Int(1),
		ReleaseType: vars.ReleaseTypeSong,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.False(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterBySince() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	since := time.Now().UTC().Truncate(time.Hour)
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Since: &since,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByTill() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	till := time.Now().UTC().Truncate(time.Hour).AddDate(0, -1, 0)
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Till: &till,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterBy_SinceAndTill() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	since := time.Now().UTC().Truncate(time.Hour).AddDate(-2, 0, 0)
	till := time.Now().UTC().Truncate(time.Hour).AddDate(0, -2, 0)
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Since: &since,
		Till:  &till,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterBy_Explicit() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	explicit := false
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Explicit: &explicit,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	// releases are sort by release date desc
	assert.NotEmpty(t.T(), rels[0].ID)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.False(t.T(), rels[0].IsExplicit)
}

func (t *testDBSuite) TestInternalReleases_Get_ByArtist_Empty() {
	// action
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{ID: 666}))
	releases, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{ArtistID: ptr.Int(666)})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}
func (t *testDBSuite) TestInternalReleases_Get_ByArtist_ArtistNotFound() {
	// action
	releases, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{ArtistID: ptr.Int(666)})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}
