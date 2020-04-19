package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
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

func getArtistIdsFromSubscriptions(subscriptions []*db.Subscription) []int64 {
	ids := make([]int64, len(subscriptions))
	for i, sub := range subscriptions {
		ids[i] = sub.ArtistID
	}
	return ids
}

func (c *SubscriptionsController) createSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	subscriptions := []*db.Subscription{}
	err = json.NewDecoder(r.Body).Decode(&subscriptions)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if len(subscriptions) == 0 {
		httputils.WriteError(w, errors.New("subscriptions weren't provided"))
		return
	}

	ids := getArtistIdsFromSubscriptions(subscriptions)
	err = db.Mgr.SubscribeUser(userName, ids)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *SubscriptionsController) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	subs, err := db.Mgr.GetUserSubscriptions(userName)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &subs)
}

func (c *SubscriptionsController) deleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	subscriptions := []*db.Subscription{}
	err = json.NewDecoder(r.Body).Decode(&subscriptions)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if len(subscriptions) == 0 {
		httputils.WriteError(w, errors.New("subscriptions weren't provided"))
		return
	}

	ids := getArtistIdsFromSubscriptions(subscriptions)
	err = db.Mgr.UnSubscribeUser(userName, ids)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
