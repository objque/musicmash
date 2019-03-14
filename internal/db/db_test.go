package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Ping_OK(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.Ping()

	// assert
	assert.NoError(t, err)
}

func TestDB_Ping_Error(t *testing.T) {
	setup()
	teardown()

	// action
	err := DbMgr.Ping()

	// assert
	assert.Error(t, err)
}
