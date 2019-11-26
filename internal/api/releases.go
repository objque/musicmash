package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
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
		httputils.WriteError(w, errors.New("since argument didn't provided"))
		return
	}

	date, err := time.Parse(dateLayout, since)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debugf("can't parse provided since argument %s as %s", since, dateLayout)
		httputils.WriteError(w, fmt.Errorf("since argument %s does not match %s layout", since, dateLayout))
		return
	}

	releases, err := db.DbMgr.FindNewReleases(date)
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &releases)
}
