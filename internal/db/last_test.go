package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_Last_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	last := time.Now().UTC()
	DbMgr.SetLastFetch(last)

	// action
	res, err := DbMgr.GetLastFetch()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, last, res.Date)
}

func TestDB_Last_Set(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.SetLastFetch(time.Now().UTC())

	// assert
	assert.NoError(t, err)
}

func TestDB_Last_Update(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	DbMgr.SetLastFetch(time.Now())

	// action
	n := time.Now().UTC()
	err := DbMgr.SetLastFetch(n)

	// assert
	assert.NoError(t, err)
	last, err := DbMgr.GetLastFetch()
	assert.Equal(t, last.Date, n)
}
