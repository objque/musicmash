package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/musicmash/internal/log"
)

func getMux() *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.DefaultCompress)

	r.Get("/healthz", healthz)

	r.Route("/v1", func(r chi.Router) {
		NewArtistsController().Register(r)
		NewAssociationsController().Register(r)
		NewSubscriptionsController().Register(r)
		NewStoresController().Register(r)
		NewNotificationSettingsController().Register(r)
		NewReleasesController().Register(r)
		NewSearchsController().Register(r)
	})
	return r
}

func ListenAndServe(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	log.Infof("Listening API on %s", addr)
	return http.ListenAndServe(addr, getMux())
}
