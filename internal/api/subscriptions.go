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
	PathSubscriptions = "/subscriptions"
)

type SubscriptionsController struct{}

func NewSubscriptionsController() *SubscriptionsController {
	return &SubscriptionsController{}
}

func (c *SubscriptionsController) Register(router chi.Router) {
	router.Route(PathSubscriptions, func(r chi.Router) {
		r.Post("/", c.createSubscriptions)
		r.Delete("/", c.deleteSubscriptions)
		r.Get("/", c.listSubscriptions)
	})
}

func (c *SubscriptionsController) createSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		WriteError(w, err)
		return
	}

	artists := []int64{}
	err = json.NewDecoder(r.Body).Decode(&artists)
	if err != nil {
		WriteError(w, errors.New("invalid body"))
		return
	}

	if len(artists) == 0 {
		WriteError(w, errors.New("artists weren't provided"))
		return
	}

	err = db.DbMgr.SubscribeUser(userName, artists)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *SubscriptionsController) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		WriteError(w, err)
		return
	}

	subs, err := db.DbMgr.GetUserSubscriptions(userName)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	b, _ := json.Marshal(&subs)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(b)
}

func (c *SubscriptionsController) deleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		WriteError(w, err)
		return
	}

	artists := []int64{}
	err = json.NewDecoder(r.Body).Decode(&artists)
	if err != nil {
		WriteError(w, errors.New("invalid body"))
		return
	}

	if len(artists) == 0 {
		WriteError(w, errors.New("artists weren't provided"))
		return
	}

	err = db.DbMgr.UnSubscribeUser(userName, artists)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
