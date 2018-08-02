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
