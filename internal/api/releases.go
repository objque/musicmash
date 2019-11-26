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
	PathReleases = "/releases"
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
	date, err := httputils.GetQueryTime(r, "since")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	releases, err := db.DbMgr.FindNewReleases(*date)
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &releases)
}
