package api

import (
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testAPISuite struct {
	suite.Suite
	server *httptest.Server
	client *api.Provider
}

func (t *testAPISuite) SetupSuite() {
	t.server = httptest.NewServer(getMux())
	t.client = api.NewProvider(t.server.URL, 1)
}

func (t *testAPISuite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.DbMgr.ApplyMigrations("../../migrations/sqlite3"))
}

func (t *testAPISuite) TearDownTest() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}

func (t *testAPISuite) TearDownSuite() {
	t.server.Close()
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(testAPISuite))
}
