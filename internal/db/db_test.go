package db

import (
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestPing_OK() {
	// action
	err := DbMgr.Ping()

	// assert
	assert.NoError(t.T(), err)
}

func (t *testDBSuite) TestPing_Error() {
	// arrange
	// close connection manually to get internal error
	DbMgr.Close()

	// action
	err := DbMgr.Ping()

	// assert
	assert.Error(t.T(), err)
}
