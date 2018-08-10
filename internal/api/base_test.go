package api

import (
	"io"
	"net/http"
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

func httpDelete(url string, body io.Reader) (resp *http.Response, err error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
