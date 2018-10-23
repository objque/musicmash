package db

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_Feed_GetUserFeedSince(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userName = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreName:  "itunes",
		StoreID:    "182821355",
		Released:   time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		Title:      "Pizza",
		StoreName:  "itunes",
		StoreID:    "213551828",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		Title:      "Pizza",
		StoreName:  "yandex",
		StoreID:    "1067",
		Released:   time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.SubscribeUserForArtists(userName, []string{"skrillex", "S.P.Y"}))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	feed, err := DbMgr.GetUserFeedSince(userName, since)

	// assert
	assert.NoError(t, err)
	assert.Len(t, feed.Announced, 1)
	assert.Len(t, feed.Released, 1)
	assert.Len(t, feed.Announced[0].Stores, 2)
	assert.Equal(t, "S.P.Y", feed.Announced[0].ArtistName)
	assert.Equal(t, "skrillex", feed.Released[0].ArtistName)
}

func TestDB_Feed_GroupReleases(t *testing.T) {
	// arrange
	releases := []*Release{
		{
			Title:      "Holy hell",
			StoreName:  "deezer",
			ArtistName: "Architects",
		},
		{
			Title:      "Wings",
			StoreName:  "deezer",
			ArtistName: "Dead birds",
		},
		{
			Title:      "Holy Hell",
			StoreName:  "itunes",
			ArtistName: "Architects",
		},
		{
			Title:      "Pirate Station History",
			ArtistName: "Pirate Station",
			StoreName:  "deezer",
		},
		{
			Title:      "Holy Hell",
			StoreName:  "spotify",
			ArtistName: "architects",
		},
		{
			Title:      "holy hell",
			StoreName:  "spotify",
			ArtistName: "Bring Me The Horizon",
		},
	}

	// action
	grouped := groupReleases(releases)

	// assert
	want := map[string]struct {
		StoresCount int
		Title       string
	}{
		"architects": {
			StoresCount: 3,
			Title:       "holy hell",
		},
		"dead birds": {
			StoresCount: 1,
			Title:       "wings",
		},
		"pirate station": {
			StoresCount: 1,
			Title:       "pirate station history",
		},
		"bring me the horizon": {
			StoresCount: 1,
			Title:       "holy hell",
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
	const posterURL = "http://pic.jpeg"
	releases := []*Release{
		{
			Title:      "Holy hell",
			StoreName:  "deezer",
			ArtistName: "Architects",
		},
		{
			Title:      "Holy hell",
			StoreName:  "spotify",
			ArtistName: "Architects",
		},
		{
			Title:      "Holy Hell",
			StoreName:  "itunes",
			ArtistName: "Architects",
			Poster:     posterURL,
		},
	}

	// action
	grouped := groupReleases(releases)

	// assert
	assert.Len(t, grouped, 1)
	assert.Equal(t, posterURL, grouped[0].Poster)
}
