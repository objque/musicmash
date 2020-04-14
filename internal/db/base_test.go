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
	suite.Run(t, new(testDBSuite))
}
