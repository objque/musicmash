package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testFetcherSuite struct {
	suite.Suite
	server *httptest.Server
	mux    *http.ServeMux
}

func (t *testFetcherSuite) SetupSuite() {
	t.mux = http.NewServeMux()
	t.server = httptest.NewServer(t.mux)
	config.Config = &config.AppConfig{
		Stores: config.StoresConfig{
			vars.StoreApple: {
				Fetch:        true,
				URL:          t.server.URL,
				Meta:         map[string]string{"token": "xxx"},
				FetchWorkers: 5,
			},
		},
	}
}

func (t *testFetcherSuite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	assert.NoError(t.T(), db.DbMgr.ApplyMigrations("../../migrations/sqlite3"))
}

func (t *testFetcherSuite) TearDownTest() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}

func (t *testFetcherSuite) TearDownSuite() {
	t.server.Close()
}

func TestFetcherSuite(t *testing.T) {
	suite.Run(t, new(testFetcherSuite))
}
