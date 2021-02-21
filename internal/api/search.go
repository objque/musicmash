package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathSearch = "/search"
)

type SearchsController struct{}

func NewSearchsController() *SearchsController {
	return &SearchsController{}
}

func (c *SearchsController) Register(router chi.Router) {
	router.Route(PathSearch, func(r chi.Router) {
		r.Get("/", c.doSearch)
	})
}

func (c *SearchsController) doSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		httputils.WriteError(w, errors.New("query argument is empty"))
		return
	}

	result, err := db.Mgr.Search(query)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &result)
}
