package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/validators"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	tasks "github.com/musicmash/musicmash/internal/tasks/subscriptions"
	"github.com/pkg/errors"
)

func createSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
		return
	}

	userArtists := []string{}
	if err := json.NewDecoder(r.Body).Decode(&userArtists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugf("User '%s' wanna subscriptions for %d artists", userName, len(userArtists))
	tasks.SubscribeUserForArtists(userName, userArtists)
	w.WriteHeader(http.StatusAccepted)
}

func deleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
		return
	}

	artists := []string{}
	if err := json.NewDecoder(r.Body).Decode(&artists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.DbMgr.UnsubscribeUserFromArtists(userName, artists); err != nil {
		log.Error(errors.Wrapf(err, "tried to unsubscribe user '%s' from artists '%v'", userName, artists))
	}

	w.WriteHeader(http.StatusOK)
}

func getUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
		return
	}

	subs, err := db.DbMgr.FindAllUserSubscriptions(userName)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to get subscriptions for '%s'", userName))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(&subs)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to marshal subscriptions for '%s'", userName))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(body)
}
