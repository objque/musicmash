package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/testutils/suite"
	"github.com/musicmash/musicmash/internal/testutils/vars"
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
	t.Suite.SetupTest()
}

func (t *testFetcherSuite) TearDownTest() {
	t.Suite.TearDownTest()
}

func (t *testFetcherSuite) TearDownSuite() {
	t.server.Close()
}

func TestFetcherSuite(t *testing.T) {
	suite.Run(t, new(testFetcherSuite))
}
