package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathNotifications = "/notifications/settings"
)

type NotificationSettingsController struct{}

func NewNotificationSettingsController() *NotificationSettingsController {
	return &NotificationSettingsController{}
}

func (s *NotificationSettingsController) Register(router chi.Router) {
	router.Route(PathNotifications, func(r chi.Router) {
		r.Post("/", s.addNotificationSettings)
		r.Patch("/", s.updateNotificationSettings)
		r.Get("/", s.listNotificationSettings)
	})
}

func (s *NotificationSettingsController) addNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	settings := db.NotificationSettings{}
	err = json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if settings.Service == "" {
		httputils.WriteError(w, errors.New("service didn't provided"))
		return
	}
	if settings.Data == "" {
		httputils.WriteError(w, errors.New("service data didn't provided"))
		return
	}

	dbSettings, err := db.DbMgr.FindNotificationSettingsForService(userName, settings.Service)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}
	if len(dbSettings) > 0 {
		httputils.WriteError(w, errors.New("user already has settings for this service"))
		return
	}

	settings.UserName = userName
	err = db.DbMgr.EnsureNotificationSettingsExists(&settings)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusCreated, &settings)
}

func (s *NotificationSettingsController) updateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	settings := db.NotificationSettings{}
	err = json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if settings.Service == "" {
		httputils.WriteError(w, errors.New("service didn't provided"))
		return
	}
	if settings.Data == "" {
		httputils.WriteError(w, errors.New("service data didn't provided"))
		return
	}

	settings.UserName = userName
	err = db.DbMgr.UpdateNotificationSettings(&settings)
	if err != nil {
		if err == db.ErrNotificationSettingsNotFound {
			httputils.WriteError(w, err)
			return
		}

		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &settings)
}

func (s *NotificationSettingsController) listNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	settings, err := db.DbMgr.FindNotificationSettings(userName)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &settings)
}
