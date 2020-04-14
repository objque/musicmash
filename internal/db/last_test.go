package db

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestLastAction_Get() {
	// arrange
	last := time.Now().UTC()
	assert.NoError(t.T(), Mgr.SetLastActionDate(ActionFetch, last))

	// action
	res, err := Mgr.GetLastActionDate(ActionFetch)

	// assert
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), last, res.Date)
}

func (t *testDBSuite) TestLastAction_Set() {
	// action
	err := Mgr.SetLastActionDate(ActionFetch, time.Now().UTC())

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestLastAction_Update() {
	// arrange
	assert.NoError(t.T(), Mgr.SetLastActionDate(ActionFetch, time.Now()))

	// action
	n := time.Now().UTC()
	err := Mgr.SetLastActionDate(ActionFetch, n)

	// assert
	assert.NoError(t.T(), err)
	last, err := Mgr.GetLastActionDate(ActionFetch)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), last.Date, n)
}

func (t *testDBSuite) TestLastAction_NotFound() {
	// arrange
	_, err := Mgr.GetLastActionDate(ActionFetch)

	// assert
	assert.Error(t.T(), err)
}
