package handler

import (
	"fmt"
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
	router := http.ServeMux{}

	return router
}
