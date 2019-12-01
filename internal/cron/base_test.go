package cron

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/suite"
)

type testCronSuite struct {
	suite.Suite
}

func (t *testCronSuite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func (t *testCronSuite) TearDownTest() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}

func TestCronSuite(t *testing.T) {
	suite.Run(t, new(testCronSuite))
}
