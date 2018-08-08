package subscriptions

import (
	"net/http"
	"net/http/httptest"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
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
		Store: config.Store{
			URL:    server.URL,
			Region: "us",
		},
	}
}

func teardown() {
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}
