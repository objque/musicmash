package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/api/validators"
)

func addUserChat(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if err := validators.IsUserExits(w, userID); err != nil {
		return
	}

	body := AddUserChatScheme{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.DbMgr.EnsureChatExists(&db.Chat{UserID: userID, ID: body.ChatID}); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
