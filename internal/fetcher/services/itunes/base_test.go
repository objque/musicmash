package itunes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/testutils/suite"
	"github.com/musicmash/musicmash/internal/testutils/vars"
)

type testAppleMusicClientSuite struct {
	suite.Suite
	server   *httptest.Server
	provider *itunes.Provider
	mux      *http.ServeMux
}

func (t *testAppleMusicClientSuite) SetupSuite() {
	t.mux = http.NewServeMux()
	t.server = httptest.NewServer(t.mux)
	t.provider = itunes.NewProvider(t.server.URL, vars.TokenSimple, time.Minute)
	config.Config = &config.AppConfig{
		Fetcher: config.FetcherConfig{
			Delay: 8,
		},
	}
}

func (t *testAppleMusicClientSuite) TearDownSuite() {
	t.server.Close()
}

func TestAppleMusicClientSuite(t *testing.T) {
	suite.Run(t, new(testAppleMusicClientSuite))
}
