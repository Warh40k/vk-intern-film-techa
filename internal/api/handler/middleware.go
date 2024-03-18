package handler

import (
	"context"
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
	log     *slog.Logger
	handler http.Handler
}

const (
	ROLE_ADMIN  = 2
	ROLE_CLIENT = 1
)

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	l.log.With(
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("since", time.Since(start).String())).
		Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.Path, time.Since(start)))
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(log *slog.Logger, handlerToWrap http.Handler) *Logger {
	return &Logger{log, handlerToWrap}
}

func (h *Handler) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		authHeader := r.Header.Get("Authorization")
		headSplit := strings.Split(authHeader, "Bearer ")
		if len(headSplit) == 2 {
			token = headSplit[1]
		} else {
			newErrResponse(h.log, w, http.StatusForbidden, r.Host+r.RequestURI, "Forbidden",
				"No Bearer token provided. Please, authorize first to access resource", "Forbidden")
			return
		}
		id, err := service.CheckJWT(token)

		if err != nil {
			newErrResponse(h.log, w, http.StatusForbidden, r.Host+r.RequestURI, "Forbidden",
				"Invalid JWT token. Please, sign up if necessary and acquire fresh token", "Forbidden")
			return
		}

		ctx := context.WithValue(r.Context(), "user", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value("user").(int)
		if !ok {
			newErrResponse(h.log, w, http.StatusForbidden, r.Host+r.RequestURI, "Forbidden",
				"Could not get user id", "Forbidden")
			return
		}
		user, err := h.services.GetUserById(id)
		if err != nil {
			newErrResponse(h.log, w, http.StatusForbidden, r.Host+r.RequestURI, "Forbidden",
				"Specified user not found", "Forbidden")
			return
		}
		if user.Role != ROLE_ADMIN {
			newErrResponse(h.log, w, http.StatusForbidden, r.Host+r.RequestURI, "Forbidden",
				"You have no admin permissions", "Forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}
