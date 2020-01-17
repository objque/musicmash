package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) fillRelease(release *Release) {
	assert.NoError(t.T(), DbMgr.EnsureReleaseExists(release))
}

func (t *testDBSuite) setupInternalReleases(id int64, r time.Time) {
	assert.NoError(t.T(), DbMgr.EnsureArtistExists(&Artist{ID: id, Name: vars.ArtistSkrillex}))

	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreApple, StoreID: "1000"})
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreDeezer, StoreID: "2000"})
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, StoreName: vars.StoreSpotify, StoreID: "3000"})

	r = r.AddDate(0, 0, -1)
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterSmall, Released: r, StoreName: vars.StoreApple, StoreID: "4000"})
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterSmall, Released: r, StoreName: vars.StoreDeezer, StoreID: "5000"})
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterSmall, Released: r, StoreName: vars.StoreSpotify, StoreID: "6000"})

	r = r.AddDate(0, -1, 0)
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseWildwaysTheX, Poster: vars.PosterMiddle, Released: r, StoreName: vars.StoreApple, StoreID: "7000"})

	r = r.AddDate(-1, 0, 0)
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: r, StoreName: vars.StoreDeezer, StoreID: "8000"})
	t.fillRelease(&Release{ArtistID: id, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: r, StoreName: vars.StoreSpotify, StoreID: "9000"})
}

func (t *testDBSuite) TestInternalReleases_GetArtist() {
	// arrange
	const artistID = 666
	r := time.Now().UTC().Truncate(time.Hour)
	t.setupInternalReleases(artistID, r)

	// action
	releases, err := DbMgr.GetArtistInternalReleases(artistID)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 4)

	expected := &InternalRelease{ID: 1, ArtistID: artistID, ArtistName: vars.ArtistSkrillex, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, ItunesID: "1000", DeezerID: "2000", SpotifyID: "3000"}
	assert.Equal(t.T(), expected, releases[0])

	r = r.AddDate(0, 0, -1)
	expected = &InternalRelease{ID: 4, ArtistID: artistID, ArtistName: vars.ArtistSkrillex, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterSmall, Released: r, ItunesID: "4000", DeezerID: "5000", SpotifyID: "6000"}
	assert.Equal(t.T(), expected, releases[1])

	r = r.AddDate(0, -1, 0)
	expected = &InternalRelease{ID: 7, ArtistID: artistID, ArtistName: vars.ArtistSkrillex, Title: vars.ReleaseWildwaysTheX, Poster: vars.PosterMiddle, Released: r, ItunesID: "7000"}
	assert.Equal(t.T(), expected, releases[2])

	r = r.AddDate(-1, 0, 0)
	expected = &InternalRelease{ID: 8, ArtistID: artistID, ArtistName: vars.ArtistSkrillex, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: r, DeezerID: "8000", SpotifyID: "9000"}
	assert.Equal(t.T(), expected, releases[3])
}

func (t *testDBSuite) TestInternalReleases_GetForUser() {
	// arrange
	const artistID = 666
	r := time.Now().UTC().Truncate(time.Hour)
	t.setupInternalReleases(artistID, r)
	// user should be subscribed after artists is created
	assert.NoError(t.T(), DbMgr.SubscribeUser(vars.UserObjque, []int64{artistID}))

	// action
	since := r.AddDate(0, -1, -3)
	till := r.AddDate(0, -1, 3)
	releases, err := DbMgr.GetUserInternalReleases(vars.UserObjque, &since, &till)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), releases, 1)
	r = r.AddDate(0, -1, -1)
	expected := &InternalRelease{ID: 7, ArtistID: artistID, ArtistName: vars.ArtistSkrillex, Title: vars.ReleaseWildwaysTheX, Poster: vars.PosterMiddle, Released: r, ItunesID: "7000"}
	assert.Equal(t.T(), expected, releases[0])
}
