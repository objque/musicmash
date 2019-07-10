package api

import (
	"net/http/httptest"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/pkg/api"
)

var (
	server *httptest.Server
	client *api.Provider
)

func setup() {
	server = httptest.NewServer(getMux())
	client = api.NewProvider(server.URL, 1)
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func teardown() {
	server.Close()
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}
