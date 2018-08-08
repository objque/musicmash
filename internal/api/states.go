package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
)

func getState(w http.ResponseWriter, r *http.Request) {
	stateID := chi.URLParam(r, "state_id")
	state, err := db.DbMgr.GetState(stateID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, _ := json.Marshal(&state)
	w.Write(buffer)
}
