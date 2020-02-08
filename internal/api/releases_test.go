package api

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) fillRelease(release *db.Release) {
	assert.NoError(t.T(), db.DbMgr.EnsureReleaseExists(release))
}

func (t *testAPISuite) TestReleases_Get() {
	// arrange
	r := time.Now().UTC()
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistAlgorithm, ID: vars.StoreIDQ}))
	assert.NoError(t.T(), db.DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreApple, StoreID: "1000", Explicit: true})
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreDeezer, StoreID: "2000", Explicit: true})
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreSpotify, StoreID: "3000", Explicit: true})

	since := r.AddDate(-1, 0, 0)
	till := r.AddDate(1, 0, 0)
	testCases := []*releases.GetOptions{
		// case 1: get releases with default period
		// (since, till will be calculated in the api)
		nil,
		// case 2: provide query args
		{Since: &since, Till: &till},
	}
	for _, opts := range testCases {
		// action
		feed, err := releases.For(t.client, vars.UserObjque, opts)

		// assert
		assert.NoError(t.T(), err)
		assert.Len(t.T(), feed, 1)
		expected := &releases.Release{
			ID:         1,
			ArtistID:   vars.StoreIDQ,
			ArtistName: vars.ArtistAlgorithm,
			Title:      vars.ReleaseAlgorithmFloatingIP,
			Poster:     vars.PosterSimple,
			Released:   r,
			ItunesID:   "1000",
			DeezerID:   "2000",
			SpotifyID:  "3000",
			Explicit:   true,
		}
		assert.Equal(t.T(), expected, feed[0])
	}
}

func (t *testAPISuite) TestReleases_Get_EmptyForPeriod() {
	// arrange
	r := time.Now().UTC()
	assert.NoError(t.T(), db.DbMgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistAlgorithm, ID: vars.StoreIDQ}))
	assert.NoError(t.T(), db.DbMgr.SubscribeUser(vars.UserObjque, []int64{vars.StoreIDQ}))
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreApple, StoreID: "1000", Explicit: true})
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreDeezer, StoreID: "2000", Explicit: true})
	t.fillRelease(&db.Release{ArtistID: vars.StoreIDQ, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreSpotify, StoreID: "3000", Explicit: true})

	// action
	since := r.AddDate(-1, -1, 0)
	till := r.AddDate(-1, 0, 0)
	releases, err := releases.For(t.client, vars.UserObjque, &releases.GetOptions{Since: &since, Till: &till})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 0)
}
