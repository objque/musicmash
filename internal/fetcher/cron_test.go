package fetcher

import (
	"testing"
	"time"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func setup() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Fetching: config.Fetching{
			Workers:                    10,
			CountOfSkippedHoursToFetch: 8,
		},
	}
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}

func TestFetcher_Internal_IsMustFetch_FirstRun(t *testing.T) {
	// first run means that no records in last_fetches
	setup()
	defer teardown()

	// action
	must := isMustFetch()

	// assert
	assert.True(t, must)
}

func TestFetcher_Internal_IsMustFetch_ReloadApp_AfterFetching(t *testing.T) {
	// fetch was successful and someone restart the app
	setup()
	defer teardown()

	// arrange
	db.DbMgr.SetLastActionDate(db.ActionFetch, time.Now().UTC())

	// action
	must := isMustFetch()

	// assert
	assert.False(t, must)
}

func TestFetcher_Internal_IsMustFetch_ReloadApp_AfterOldestFetching(t *testing.T) {
	// fetch was successful some times ago and someone restart the app
	setup()
	defer teardown()

	// arrange
	db.DbMgr.SetLastActionDate(db.ActionFetch, time.Now().UTC().Add(-time.Hour * 48))

	// action
	must := isMustFetch()

	// assert
	assert.True(t, must)
}
