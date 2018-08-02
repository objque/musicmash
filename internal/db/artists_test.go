package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Artists_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExists(&Artist{
		Name: "skrillex",
	})

	// assert
	assert.NoError(t, err)
	artist, err := DbMgr.FindArtistByName("skrillex")
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", artist.Name)
}

func TestDB_Artists_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{
		Name: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{
		Name: "S.P.Y",
	}))

	// action
	artists, err := DbMgr.GetAllArtists()

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, "skrillex", artists[0].Name)
	assert.Equal(t, "S.P.Y", artists[1].Name)
}
