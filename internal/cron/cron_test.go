package cron

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func (t *testCronSuite) TestIsMustFetch_FirstRun() {
	// first run means that no records in last_fetches
	c := cron{ActionName: db.ActionFetch}

	// action
	must := c.IsMustFetch()

	// assert
	assert.False(t.T(), must)
}

func (t *testCronSuite) TestIsMustFetch_ReloadApp_AfterFetching() {
	// fetch was successful and someone restart the app
	c := cron{ActionName: db.ActionFetch, CountOfSkippedHoursToRun: 8}

	// arrange
	assert.NoError(t.T(), db.DbMgr.SetLastActionDate(db.ActionFetch, time.Now().UTC()))

	// action
	must := c.IsMustFetch()

	// assert
	assert.False(t.T(), must)
}

func (t *testCronSuite) TestIsMustFetch_ReloadApp_AfterOldestFetching() {
	// fetch was successful some times ago and someone restart the app
	c := cron{ActionName: db.ActionFetch}

	// arrange
	assert.NoError(t.T(), db.DbMgr.SetLastActionDate(db.ActionFetch, time.Now().UTC().Add(-time.Hour*48)))

	// action
	must := c.IsMustFetch()

	// assert
	assert.True(t.T(), must)
}
