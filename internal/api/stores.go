package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
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
		r.Post("/", s.addStore)
		r.Get("/", s.listStores)
	})
}

func (s *StoresController) addStore(w http.ResponseWriter, r *http.Request) {
	store := db.Store{}
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		WriteError(w, errors.New("invalid body"))
		return
	}

	if store.Name == "" {
		WriteError(w, errors.New("store name didn't provided"))
		return
	}

	if exists := db.DbMgr.IsStoreExists(store.Name); exists {
		WriteError(w, errors.New("store already exists"))
		return
	}

	err = db.DbMgr.EnsureStoreExists(store.Name)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	b, _ := json.Marshal(&store)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(b)
}

func (s *StoresController) listStores(w http.ResponseWriter, _ *http.Request) {
	stores, err := db.DbMgr.GetAllStores()
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	buffer, err := json.Marshal(&stores)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(buffer)
}
