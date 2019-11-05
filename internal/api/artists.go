package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathArtists = "/artists"
)

type ArtistsController struct{}

func NewArtistsController() *ArtistsController {
	return &ArtistsController{}
}

func (c *ArtistsController) Register(router chi.Router) {
	router.Route(PathArtists, func(r chi.Router) {
		r.Post("/", c.addArtist)
		r.Post("/associate", c.associateArtist)
		r.Get("/", c.searchArtist)
	})
}

func (c *ArtistsController) addArtist(w http.ResponseWriter, r *http.Request) {
	artist := db.Artist{}
	err := json.NewDecoder(r.Body).Decode(&artist)
	if err != nil {
		WriteError(w, errors.New("invalid body"))
		return
	}

	// do not allow override ID
	if artist.Name == "" {
		WriteError(w, errors.New("artist name didn't provided"))
		return
	}
	artist.ID = 0

	err = db.DbMgr.EnsureArtistExists(&artist)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	b, _ := json.Marshal(&artist)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(b)
}

func (c *ArtistsController) associateArtist(w http.ResponseWriter, r *http.Request) {
	info := db.ArtistStoreInfo{}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		WriteError(w, errors.New("invalid body"))
		return
	}

	if exist := db.DbMgr.IsArtistExistsInStore(info.StoreID, info.StoreName); exist {
		WriteError(w, errors.New("artist already associated"))
		return
	}

	_, err = db.DbMgr.GetArtistWithFullInfo(info.ArtistID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			WriteError(w, errors.New("artist not found"))
			return
		}

		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	err = db.DbMgr.EnsureArtistExistsInStore(info.ArtistID, info.StoreName, info.StoreID)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	b, _ := json.Marshal(&info)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(b)
}

func (c *ArtistsController) searchArtist(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		WriteError(w, errors.New("name didn't provided"))
		return
	}

	artists, err := db.DbMgr.SearchArtists(name)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	buffer, err := json.Marshal(&artists)
	if err != nil {
		WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(buffer)
}
