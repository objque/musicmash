package db

import (
	"time"

	"github.com/stretchr/testify/assert"
)

const ActionFetch = "fetch"

func (t *testDBSuite) TestLastAction_Get() {
	// arrange
	last := time.Now().UTC().Truncate(time.Minute)
	assert.NoError(t.T(), Mgr.SetLastActionDate(ActionFetch, last))

	// action
	res, err := Mgr.GetLastActionDate(ActionFetch)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), last, res.Date.UTC())
}

func (t *testDBSuite) TestLastAction_Set() {
	// action
	now := time.Now().UTC().Truncate(time.Minute)
	err := Mgr.SetLastActionDate(ActionFetch, now)

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestLastAction_Update() {
	// arrange
	assert.NoError(t.T(), Mgr.SetLastActionDate(ActionFetch, time.Now()))

	// action
	now := time.Now().UTC().Truncate(time.Minute)
	err := Mgr.SetLastActionDate(ActionFetch, now)

	// assert
	assert.NoError(t.T(), err)
	last, err := Mgr.GetLastActionDate(ActionFetch)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), last.Date.UTC(), now)
}

func (t *testDBSuite) TestLastAction_NotFound() {
	// arrange
	_, err := Mgr.GetLastActionDate(ActionFetch)

	// assert
	assert.Error(t.T(), err)
}
