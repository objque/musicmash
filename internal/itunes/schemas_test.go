package itunes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelease_IsLatest(t *testing.T) {
	release := Release{
		ReleaseDate: time.Now().UTC(),
	}

	assert.True(t, release.IsLatest())
}

func TestRelease_IsLatest_False(t *testing.T) {
	release := Release{
		ReleaseDate: time.Now().UTC().Add(-time.Hour * 700),
	}

	assert.False(t, release.IsLatest())
}
