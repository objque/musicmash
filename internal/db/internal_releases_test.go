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
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistAlgorithm, ID: vars.StoreIDQ}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex, ID: vars.StoreIDW}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))
	// first release
	t.fillRelease(&Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreApple, StoreID: "1000", Type: vars.ReleaseTypeAlbum, Explicit: true})
	t.fillRelease(&Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreDeezer, StoreID: "2000", Type: vars.ReleaseTypeAlbum, Explicit: true})
	t.fillRelease(&Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreSpotify, StoreID: "3000", Type: vars.ReleaseTypeAlbum, Explicit: true})
	// second release
	t.fillRelease(&Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterMiddle, Released: monthAgo, StoreName: vars.StoreApple, StoreID: "1100", Type: vars.ReleaseTypeSong, Explicit: false})
	// third release from another artist
	t.fillRelease(&Release{ArtistID: vars.StoreIDW, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: yearAgo, StoreName: vars.StoreApple, StoreID: "1110", Type: vars.ReleaseTypeVideo, Explicit: true})
}

func (t *testDBSuite) TestInternalReleases_Get_All() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(nil)
	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 3)
	// releases are sort by release date desc
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "1000", *rels[0].ItunesID)
	assert.Equal(t.T(), "2000", *rels[0].DeezerID)
	assert.Equal(t.T(), "3000", *rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)

	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", *rels[1].ItunesID)
	assert.Nil(t.T(), rels[1].SpotifyID)
	assert.Nil(t.T(), rels[1].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.False(t.T(), rels[1].Explicit)

	assert.Equal(t.T(), int64(vars.StoreIDW), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[2].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[2].Poster)
	assert.Equal(t.T(), "1110", *rels[2].ItunesID)
	assert.Nil(t.T(), rels[2].SpotifyID)
	assert.Nil(t.T(), rels[2].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[2].Type)
	assert.True(t.T(), rels[2].Explicit)
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
	assert.Equal(t.T(), int64(vars.StoreIDW), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)

	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", *rels[1].ItunesID)
	assert.Nil(t.T(), rels[1].SpotifyID)
	assert.Nil(t.T(), rels[1].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.False(t.T(), rels[1].Explicit)

	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[2].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[2].Poster)
	assert.Equal(t.T(), "1000", *rels[2].ItunesID)
	assert.Equal(t.T(), "2000", *rels[2].DeezerID)
	assert.Equal(t.T(), "3000", *rels[2].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[2].Type)
	assert.True(t.T(), rels[2].Explicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByLimit() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Limit: ptr.Uint(1),
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "1000", *rels[0].ItunesID)
	assert.Equal(t.T(), "2000", *rels[0].DeezerID)
	assert.Equal(t.T(), "3000", *rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterBy_LimitAndOffset() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		Limit:  ptr.Uint(1),
		Offset: ptr.Uint(1),
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.False(t.T(), rels[0].Explicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByArtistID() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		ArtistID: ptr.Int(vars.StoreIDQ),
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 2)
	// releases are sort by release date desc
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "1000", *rels[0].ItunesID)
	assert.Equal(t.T(), "2000", *rels[0].DeezerID)
	assert.Equal(t.T(), "3000", *rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)

	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", *rels[1].ItunesID)
	assert.Nil(t.T(), rels[1].SpotifyID)
	assert.Nil(t.T(), rels[1].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.False(t.T(), rels[1].Explicit)
}

func (t *testDBSuite) TestInternalReleases_Get_FilterByArtistID_ReleaseType() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := Mgr.GetInternalReleases(&GetInternalReleaseOpts{
		ArtistID:    ptr.Int(vars.StoreIDQ),
		ReleaseType: vars.ReleaseTypeSong,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.False(t.T(), rels[0].Explicit)
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
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "1000", *rels[0].ItunesID)
	assert.Equal(t.T(), "2000", *rels[0].DeezerID)
	assert.Equal(t.T(), "3000", *rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)
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
	assert.Equal(t.T(), int64(vars.StoreIDW), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)
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
	assert.Equal(t.T(), int64(vars.StoreIDW), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.True(t.T(), rels[0].Explicit)
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
	assert.Equal(t.T(), int64(vars.StoreIDQ), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", *rels[0].ItunesID)
	assert.Nil(t.T(), rels[0].SpotifyID)
	assert.Nil(t.T(), rels[0].DeezerID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.False(t.T(), rels[0].Explicit)
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
