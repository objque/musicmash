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

func (t *Suite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.DbMgr.ApplyMigrations(db.GetPathToMigrationsDir()))
}

func (t *Suite) TearDownTest() {
	_ = db.DbMgr.Close()
}

func Run(t *testing.T, testingSuite suite.TestingSuite) {
	suite.Run(t, testingSuite)
}
