package handler

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"log/slog"
	"net/http"
)

const (
	sortRating   = "rating"
	sortTitle    = "title"
	sortReleased = "released"
	ascSort      = "asc"
	descSort     = "desc"
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
	router.HandleFunc("POST /signup/", h.SignUp)
	router.HandleFunc("POST /auth/", h.SignIn)

	router.Handle("POST /films/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.CreateFilm))))
	router.Handle("GET /films/", h.CheckAuth(http.HandlerFunc(h.ListFilms)))
	router.Handle("GET /films/search/", h.CheckAuth(http.HandlerFunc(h.SearchFilm)))

	router.Handle("PUT /films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.UpdateFilm))))
	router.Handle("PATCH /films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.PatchFilm))))
	router.Handle("DELETE /films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.DeleteFilm))))

	router.Handle("GET /actors/", h.CheckAuth(http.HandlerFunc(h.ListActors)))
	router.Handle("POST /actors/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.CreateActor))))

	router.Handle("PUT /actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.UpdateActor))))
	router.Handle("PATCH /actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.PatchActor))))
	router.Handle("DELETE /actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.DeleteActor))))

	return router
}
