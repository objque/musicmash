package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathReleases = "/releases"
	dateLayout   = "2006-01-02T15:04:05"
)

type ReleasesController struct{}

func NewReleasesController() *ReleasesController {
	return &ReleasesController{}
}

func (c *ReleasesController) Register(router chi.Router) {
	router.Route(PathReleases, func(r chi.Router) {
		r.Get("/", c.getReleases)
	})
}

func (c *ReleasesController) getReleases(w http.ResponseWriter, r *http.Request) {
	since := r.URL.Query().Get("since")
	if len(since) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateLayout, since)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debugf("can't parse provided since argument '%s' as '%s'", since, dateLayout)
		return
	}

	releases, err := db.DbMgr.FindNewReleases(date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	buffer, err := json.Marshal(&releases)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(buffer)
}
