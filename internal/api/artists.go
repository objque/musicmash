package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/api/validators"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

func searchArtist(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	if err := validators.IsUserExits(w, userName); err != nil {
		return
	}

	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	artists, err := db.DbMgr.SearchArtists(name)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, _ := json.Marshal(&artists)
	_, _ = w.Write(buffer)
}

func getArtistDetails(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(chi.URLParam(r, "artist_name"))
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	details, err := db.DbMgr.GetArtistDetails(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, _ := json.Marshal(&details)
	_, _ = w.Write(buffer)
}
