package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testDBSuite struct {
	suite.Suite
}

func (t *testDBSuite) SetupTest() {
	DbMgr = NewFakeDatabaseMgr()
	assert.NoError(t.T(), DbMgr.ApplyMigrations("../../migrations/sqlite3"))
}

func (t *testDBSuite) TearDownTest() {
	_ = DbMgr.DropAllTables()
	_ = DbMgr.Close()
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(testDBSuite))
}
