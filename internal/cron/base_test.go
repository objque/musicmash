package cron

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testCronSuite struct {
	suite.Suite
}

func (t *testCronSuite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.DbMgr.ApplyMigrations("../../migrations/sqlite3"))
}

func (t *testCronSuite) TearDownTest() {
	_ = db.DbMgr.Close()
}

func TestCronSuite(t *testing.T) {
	suite.Run(t, new(testCronSuite))
}
