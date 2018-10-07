package api

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/db"
)

var (
	server *httptest.Server
	mux    *chi.Mux
)

func setup() {
	db.DbMgr = db.NewFakeDatabaseMgr()
	server = httptest.NewServer(getMux())
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
