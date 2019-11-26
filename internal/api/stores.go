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
	PathStores = "/stores"
)

type StoresController struct{}

func NewStoresController() *StoresController {
	return &StoresController{}
}

func (s *StoresController) Register(router chi.Router) {
	router.Route(PathStores, func(r chi.Router) {
		r.Get("/", s.listStores)
	})
}

func (s *StoresController) listStores(w http.ResponseWriter, _ *http.Request) {
	stores, err := db.DbMgr.GetAllStores()
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &stores)
}
