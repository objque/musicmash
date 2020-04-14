package db

import (
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestPing_OK() {
	// action
	err := Mgr.Ping()

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestPing_Error() {
	// arrange
	// close connection manually to get internal error
	assert.NoError(t.T(), Mgr.Close())

	// action
	err := Mgr.Ping()

	// assert
	assert.Error(t.T(), err)
}
