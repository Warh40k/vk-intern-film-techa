package app

import (
	"context"
	"net/http"
	"time"
)

type App struct {
	httpServer *http.Server
}

func (s *App) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *App) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
