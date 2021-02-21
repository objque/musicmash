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
