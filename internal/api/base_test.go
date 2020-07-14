package api

import (
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/suite"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/stretchr/testify/assert"
)

type testAPISuite struct {
	suite.Suite
	server *httptest.Server
	client *api.Provider
}

func (t *testAPISuite) SetupSuite() {
	t.server = httptest.NewServer(getMux())
	t.client = api.NewProvider(t.server.URL, 1)
	db.Mgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.Mgr.ApplyMigrations(db.GetPathToMigrationsDir()))
}

func (t *testAPISuite) TearDownTest() {
	assert.NoError(t.T(), db.Mgr.TruncateAllTables())
}

func (t *testAPISuite) TearDownSuite() {
	assert.NoError(t.T(), db.Mgr.DropAllTablesViaMigrations(db.GetPathToMigrationsDir()))
	_ = db.Mgr.Close()
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping generic suite tests in short mode")
	}

	suite.Run(t, new(testAPISuite))
}
