package api

import (
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/stretchr/testify/suite"
)

type testApiSuite struct {
	suite.Suite
	server *httptest.Server
	client *api.Provider
}

func (t *testApiSuite) SetupSuite() {
	t.server = httptest.NewServer(getMux())
	t.client = api.NewProvider(t.server.URL, 1)
}

func (t *testApiSuite) SetupTest() {
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func (t *testApiSuite) TearDownTest() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}

func (t *testApiSuite) TearDownSuite() {
	t.server.Close()
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(testApiSuite))
}
