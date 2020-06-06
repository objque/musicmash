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
	Mgr = NewFakeDatabaseMgr()
	assert.NoError(t.T(), Mgr.ApplyMigrations(GetPathToMigrationsDir()))
}

func (t *testDBSuite) TearDownTest() {
	_ = Mgr.Close()
}

func TestDBSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db tests in short mode")
	}

	suite.Run(t, new(testDBSuite))
}
