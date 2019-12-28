package api

import (
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/testutils/suite"
	"github.com/musicmash/musicmash/pkg/api"
)

type testAPISuite struct {
	suite.Suite
	server *httptest.Server
	client *api.Provider
}

func (t *testAPISuite) SetupSuite() {
	t.Suite.SetupTest()
	t.server = httptest.NewServer(getMux())
	t.client = api.NewProvider(t.server.URL, 1)
}

func (t *testAPISuite) SetupTest() {
	t.Suite.SetupTest()
}

func (t *testAPISuite) TearDownTest() {
	t.Suite.TearDownTest()
}

func (t *testAPISuite) TearDownSuite() {
	t.server.Close()
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(testAPISuite))
}
