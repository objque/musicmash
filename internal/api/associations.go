package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathAssociations = "/associations"
)

type AssociationsController struct{}

func NewAssociationsController() *AssociationsController {
	return &AssociationsController{}
}

func (c *AssociationsController) Register(router chi.Router) {
	router.Route(PathAssociations, func(r chi.Router) {
		r.Post("/", c.addAssociation)
		r.Get("/", c.listAssociations)
	})
}

func (c *AssociationsController) addAssociation(w http.ResponseWriter, r *http.Request) {
	info := db.Association{}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if exist := db.Mgr.IsStoreExists(info.StoreName); !exist {
		httputils.WriteError(w, errors.New("store not found"))
		return
	}

	if exist := db.Mgr.IsAssociationExists(info.StoreName, info.StoreID); exist {
		httputils.WriteError(w, errors.New("artist already associated"))
		return
	}

	_, err = db.Mgr.GetArtist(info.ArtistID)
	if err != nil {
		if err == sql.ErrNoRows {
			httputils.WriteError(w, errors.New("artist not found"))
			return
		}

		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	err = db.Mgr.EnsureAssociationExists(info.ArtistID, info.StoreName, info.StoreID)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusCreated, &info)
}

//nolint:gocognit
func (c *AssociationsController) listAssociations(w http.ResponseWriter, r *http.Request) {
	opts := db.AssociationOpts{}

	if storeName := r.URL.Query().Get("store_name"); storeName != "" {
		if !db.Mgr.IsStoreExists(storeName) {
			httputils.WriteError(w, errors.New("store_name not exists"))
			return
		}

		opts.StoreName = storeName
	}

	if id := r.URL.Query().Get("artist_id"); id != "" {
		artistID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			httputils.WriteError(w, errors.New("artist_id should be int"))
			return
		}

		_, err = db.Mgr.GetArtist(artistID)
		if err != nil {
			if err == sql.ErrNoRows {
				httputils.WriteError(w, errors.New("artist not found"))
				return
			}

			httputils.WriteInternalError(w)
			log.Error(err)
			return
		}

		opts.ArtistID = artistID
	}

	associations, err := db.Mgr.FindAssociations(&opts)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &associations)
}
