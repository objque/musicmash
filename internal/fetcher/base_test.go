package fetcher

import (
	"net/http"
	"net/http/httptest"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	db.DbMgr = db.NewFakeDatabaseMgr()
	config.Config = &config.AppConfig{
		Stores: []*config.Store{
			{Name: testutil.StoreYandex, URL: server.URL, FetchWorkers: 1},
			{Name: testutil.StoreApple, URL: server.URL, Meta: map[string]string{"token": "xxx"}, FetchWorkers: 1},
		},
	}
}

func teardown() {
	server.Close()
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}
