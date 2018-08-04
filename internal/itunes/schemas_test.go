package itunes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLastRelease_IsLatest(t *testing.T) {
	release := LastRelease{
		Date: time.Now().UTC(),
	}

	assert.True(t, release.IsLatest())
}

func TestLastRelease_IsLatest_False(t *testing.T) {
	release := LastRelease{
		Date: time.Now().UTC().Add(-time.Hour * 700),
	}

	assert.False(t, release.IsLatest())
}

func TestRelease_GetCollectionType(t *testing.T) {
	// arrange
	releases := []struct {
		TracksCount int
		Title       string
		Type        string
	}{
		{
			Type:        EPReleaseType,
			Title:       "The Remedy",
			TracksCount: 4,
		},
		{
			Type:        SingleReleaseType,
			Title:       "City Lights - Single",
			TracksCount: 1,
		},
		{
			Type:        SingleReleaseType,
			Title:       "Secrets (The Remixes) - Single",
			TracksCount: 3,
		},
		{
			Type:        AlbumReleaseType,
			Title:       "Only the Best",
			TracksCount: 37,
		},
		{
			Type:        AlbumReleaseType,
			Title:       "KIDS SEE GHOSTS",
			TracksCount: 7,
		},
		{
			Type:        EPReleaseType,
			Title:       "Thor Ep",
			TracksCount: 7,
		},
	}

	for _, release := range releases {
		result := Release{TrackCount: release.TracksCount, CollectionName: release.Title}
		// action and assert
		assert.Equal(t, release.Type, result.GetCollectionType())
	}
}
