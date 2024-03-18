package handler

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://locahost:8080/swagger/swagger.json")))

	router.HandleFunc("POST /api/v1/signup/", h.SignUp)
	router.HandleFunc("POST /api/v1/auth/", h.SignIn)

	router.Handle("POST /api/v1/films/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.CreateFilm))))
	router.Handle("GET /api/v1/films/", h.CheckAuth(http.HandlerFunc(h.ListFilms)))
	router.Handle("GET /api/v1/films/search/", h.CheckAuth(http.HandlerFunc(h.SearchFilm)))

	router.Handle("PUT /api/v1/films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.UpdateFilm))))
	router.Handle("PATCH /api/v1/films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.PatchFilm))))
	router.Handle("DELETE /api/v1/films/{film_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.DeleteFilm))))

	router.Handle("GET /api/v1/actors/", h.CheckAuth(http.HandlerFunc(h.ListActors)))
	router.Handle("POST /api/v1/actors/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.CreateActor))))

	router.Handle("PUT /api/v1/actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.UpdateActor))))
	router.Handle("PATCH /api/v1/actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.PatchActor))))
	router.Handle("DELETE /api/v1/actors/{actor_id}/", h.CheckAuth(h.CheckAdmin(http.HandlerFunc(h.DeleteActor))))

	return router
}
