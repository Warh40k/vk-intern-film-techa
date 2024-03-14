package handler

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"log/slog"
	"net/http"
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

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /films/", h.CreateFilm)
	router.HandleFunc("GET /films/", h.ListFilms)
	router.HandleFunc("GET /films/search/", h.SearchFilm)

	router.HandleFunc("PUT /films/{film_id}/", h.UpdateFilm)
	router.HandleFunc("PATCH /films/{film_id}/", h.PatchFilm)
	router.HandleFunc("DELETE /films/{film_id}/", h.DeleteFilm)

	router.HandleFunc("GET /actors/", h.ListActors)
	router.HandleFunc("POST /actors/", h.CreateActor)

	router.HandleFunc("PUT /actors/{actor_id}/", h.UpdateActor)
	router.HandleFunc("PATCH /actors/{actor_id}/", h.PatchActor)
	router.HandleFunc("DELETE /actors/{actor_id}/", h.DeleteActor)

	return router
}
