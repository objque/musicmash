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
		Tasks: config.Tasks{
			Subscriptions: config.SubscriptionsTask{
				FindArtistWorkers:      1,
				SubscribeArtistWorkers: 1,
			},
		},
	}
}

func teardown() {
	//db.DbMgr.DropAllTables()
	db.DbMgr.Close()
}
