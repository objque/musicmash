package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testDBSuite struct {
	suite.Suite
}

func (t *testDBSuite) SetupSuite() {
	Mgr = NewFakeDatabaseMgr()
	assert.NoError(t.T(), Mgr.ApplyMigrations(fmt.Sprintf("file://%s", GetPathToMigrationsDir())))
}

func (t *testDBSuite) TearDownTest() {
	assert.NoError(t.T(), Mgr.TruncateAllTables())
}

func (t *testDBSuite) TearDownSuite() {
	assert.NoError(t.T(), Mgr.DropAllTablesViaMigrations(fmt.Sprintf("file://%s", GetPathToMigrationsDir())))
	_ = Mgr.Close()
}

func TestDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db tests in short mode")
	}

	suite.Run(t, new(testDBSuite))
}
