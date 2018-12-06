package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/validators"
	"github.com/musicmash/musicmash/internal/feed"
	"github.com/musicmash/musicmash/internal/log"
)

func getUserFeed(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
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

	userFeed, err := feed.GetForUser(userName, weekAgo, time.Now().UTC())
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("format")
	var body []byte
	switch strings.ToLower(format) {
	case "rss":
		body, err = feed.Formatter.ToRss(userFeed)
	default:
		body, err = feed.Formatter.ToJson(userFeed)
	}
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
