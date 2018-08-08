package api

import (
	"net/http/httptest"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
)

var server *httptest.Server

func setup() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	server = httptest.NewServer(getMux())
	config.Config = &config.AppConfig{
		Store: config.Store{
			URL:    server.URL,
			Region: "us",
		},
		Tasks: config.Tasks{
			Subscriptions: config.SubscriptionsTask{
				UseSearchDelay:         false,
				SubscribeArtistWorkers: 10,
			},
		},
	}
}

func teardown() {
	server.Close()
	db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}
