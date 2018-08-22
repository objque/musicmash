package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_LastAction_Get(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	last := time.Now().UTC()
	DbMgr.SetLastActionDate(ActionFetch, last)

	// action
	res, err := DbMgr.GetLastActionDate(ActionFetch)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, last, res.Date)
}

func TestDB_LastAction_Set(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.SetLastActionDate(ActionFetch, time.Now().UTC())

	// assert
	assert.NoError(t, err)
}

func TestDB_LastAction_Update(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	DbMgr.SetLastActionDate(ActionFetch, time.Now())

	// action
	n := time.Now().UTC()
	err := DbMgr.SetLastActionDate(ActionFetch, n)

	// assert
	assert.NoError(t, err)
	last, err := DbMgr.GetLastActionDate(ActionFetch)
	assert.Equal(t, last.Date, n)
}

func TestDB_LastAction_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	_, err := DbMgr.GetLastActionDate(ActionFetch)

	// assert
	assert.Error(t, err)
}
