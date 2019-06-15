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
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
		return
	}

	since := r.URL.Query().Get("since")
	weekAgo := time.Now().UTC().Add(-time.Hour * 24 * 7).Truncate(time.Hour * 24)
	if since != "" {
		var err error
		weekAgo, err = time.Parse("2006-01-02", since)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	feed, err := db.DbMgr.GetUserFeedSince(userName, weekAgo)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, _ := json.Marshal(&feed)
	_, _ = w.Write(buffer)
}
