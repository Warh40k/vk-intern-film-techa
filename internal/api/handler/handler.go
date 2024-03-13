package handler

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

type Handler struct {
	services *service.Service
	log      *slog.Logger
}

func NewHandler(services *service.Service, log *slog.Logger) *Handler {
	return &Handler{services: services, log: log}
}

type UploadError struct {
	Filename string
	Message  string
	Err      error
}

func (e UploadError) Error() string {
	return fmt.Sprintf("%s: %s", e.Filename, e.Message)
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Route("/films", func(r chi.Router) {
		r.Get("/", h.ListFilms)
		r.Post("/", h.CreateFilm)
		r.Get("/search", h.SearchFilm)
		r.Route("/{film_id}", func(r chi.Router) {
			r.Put("/", h.UpdateFilm)
			r.Patch("/", h.PatchFilm)
			r.Delete("/", h.DeleteFilm)
		})
	})
	router.Route("/actors", func(r chi.Router) {
		r.Get("/", h.ListActors)
		r.Post("/", h.CreateActor)
		r.Route("/{actor_id}", func(r chi.Router) {
			r.Patch("/", h.PatchActor)
			r.Put("/", h.UpdateActor)
			r.Delete("/", h.DeleteActor)
		})
	})

	return router
}
