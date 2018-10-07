package fetcher

import (
	"net/http"
	"net/http/httptest"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
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
			{Name: "yandex", URL: server.URL},
			{Name: "itunes", URL: server.URL, Meta: map[string]string{"token": "xxx"}},
		},
		Fetching: config.Fetching{
			Workers: 1,
		},
	}
}

func teardown() {
	server.Close()
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}
