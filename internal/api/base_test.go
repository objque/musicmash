package api

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
)

var (
	server *httptest.Server
)

func setup() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	server = httptest.NewServer(getMux())
	config.Config = &config.AppConfig{
		Stores: map[string]*config.Store{
			testutil.StoreApple: {ReleaseURL: "https://itunes.apple.com/us/album/%s"},
		},
	}
}

func teardown() {
	server.Close()
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
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
