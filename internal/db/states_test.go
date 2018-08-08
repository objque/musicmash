package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_States_Update(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.UpdateState("xxx", ProcessingState)

	// assert
	assert.NoError(t, err)
	state, err := DbMgr.GetState("xxx")
	assert.NoError(t, err)
	assert.Equal(t, ProcessingState, state.Value)
}
