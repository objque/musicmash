package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/validators"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

func getUserFeed(w http.ResponseWriter, r *http.Request) {
	UserName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, UserName); err != nil {
		return
	}

	since := r.URL.Query().Get("since")
	weekAgo := time.Now().UTC().Add(-time.Hour * 24 * 7)
	if since != "" {
		var err error
		weekAgo, err = time.Parse("2006-01-02", since)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	feed, err := db.DbMgr.GetUserFeedSince(UserName, weekAgo)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, _ := json.Marshal(&feed)
	w.Write(buffer)
}
