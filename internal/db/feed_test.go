package db

import (
	"strings"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Feed_GetUserFeedSince(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSkrillex,
		StoreName:  testutil.StoreApple,
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		Title:      "Pizza",
		StoreName:  testutil.StoreApple,
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: testutil.ArtistSPY,
		Title:      "Pizza",
		StoreName:  testutil.StoreYandex,
		StoreID:    "1067",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(testutil.UserObjque, []string{testutil.ArtistSkrillex, testutil.ArtistSPY}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	feed, err := DbMgr.GetUserFeedSince(testutil.UserObjque, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, feed.Announced, 1)
	assert.Len(t, feed.Released, 1)
	assert.Len(t, feed.Announced[0].Stores, 2)
	assert.Equal(t, testutil.ArtistSPY, feed.Announced[0].ArtistName)
	assert.Equal(t, testutil.ArtistSkrillex, feed.Released[0].ArtistName)
}

func TestDB_Feed_GroupReleases(t *testing.T) {
	// arrange
	releases := []*Release{
		{
			Title:      strings.ToLower(testutil.ReleaseArchitectsHollyHell),
			StoreName:  testutil.StoreDeezer,
			ArtistName: testutil.ArtistArchitects,
		},
		{
			Title:      testutil.ReleaseWildwaysTheX,
			StoreName:  testutil.StoreDeezer,
			ArtistName: testutil.ArtistWildways,
		},
		{
			Title:      testutil.ReleaseArchitectsHollyHell,
			StoreName:  testutil.StoreApple,
			ArtistName: testutil.ArtistArchitects,
		},
		{
			Title:      testutil.ReleaseSkrillexRecess,
			ArtistName: testutil.ArtistSkrillex,
			StoreName:  testutil.StoreDeezer,
		},
		{
			Title:      testutil.ReleaseArchitectsHollyHell,
			StoreName:  testutil.StoreSpotify,
			ArtistName: strings.ToLower(testutil.ArtistArchitects),
		},
		{
			Title:      testutil.ReleaseAlgorithmFloatingIP,
			StoreName:  testutil.StoreSpotify,
			ArtistName: testutil.ArtistAlgorithm,
		},
	}

	// action
	grouped := groupReleases(releases)

	// assert
	want := map[string]struct {
		StoresCount int
		Title       string
	}{
		strings.ToLower(testutil.ArtistArchitects): {
			StoresCount: 3,
			Title:       strings.ToLower(testutil.ReleaseArchitectsHollyHell),
		},
		strings.ToLower(testutil.ArtistWildways): {
			StoresCount: 1,
			Title:       strings.ToLower(testutil.ReleaseWildwaysTheX),
		},
		strings.ToLower(testutil.ArtistSkrillex): {
			StoresCount: 1,
			Title:       strings.ToLower(testutil.ReleaseSkrillexRecess),
		},
		strings.ToLower(testutil.ArtistAlgorithm): {
			StoresCount: 1,
			Title:       strings.ToLower(testutil.ReleaseAlgorithmFloatingIP),
		},
	}
	assert.Len(t, grouped, 4)
	for _, release := range grouped {
		val, ok := want[strings.ToLower(release.ArtistName)]
		assert.True(t, ok)
		assert.Len(t, release.Stores, val.StoresCount)
		assert.Equal(t, val.Title, strings.ToLower(release.Title))
	}
}

func TestDB_Feed_GroupReleases_OverridePoster_IfWasEmpty(t *testing.T) {
	// arrange
	releases := []*Release{
		{
			Title:      strings.ToLower(testutil.ReleaseArchitectsHollyHell),
			StoreName:  testutil.StoreDeezer,
			ArtistName: testutil.ArtistArchitects,
		},
		{
			Title:      strings.ToLower(testutil.ReleaseArchitectsHollyHell),
			StoreName:  testutil.StoreSpotify,
			ArtistName: testutil.ArtistArchitects,
		},
		{
			Title:      testutil.ReleaseArchitectsHollyHell,
			StoreName:  testutil.StoreApple,
			ArtistName: testutil.ArtistArchitects,
			Poster:     testutil.PosterSimple,
		},
	}

	// action
	grouped := groupReleases(releases)

	// assert
	assert.Len(t, grouped, 1)
	assert.Equal(t, testutil.PosterSimple, grouped[0].Poster)
}
