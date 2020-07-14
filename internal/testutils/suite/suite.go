package suite

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func (t *Suite) SetupSuite() {
	db.Mgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.Mgr.ApplyMigrations(db.GetPathToMigrationsDir()))
}

func (t *Suite) TearDownTest() {
	assert.NoError(t.T(), db.Mgr.TruncateAllTables())
}

func (t *Suite) TearDownSuite() {
	assert.NoError(t.T(), db.Mgr.DropAllTablesViaMigrations(db.GetPathToMigrationsDir()))
	_ = db.Mgr.Close()
}

func Run(t *testing.T, testingSuite suite.TestingSuite) {
	if testing.Short() {
		t.Skip("skipping generic suite tests in short mode")
	}

	suite.Run(t, testingSuite)
}
