package handler

import (
	"encoding/json"
	"errors"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required,gte=1,lte=128"`
	Password string `json:"password" validate:"required,gte=8,lte=128"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

// SignIn godoc
//
//		@Summary		Авторизация
//		@Description	Получения токена авторизации
//		@Tags			auth
//		@Accept			json
//		@Produce		json
//	 	@Param			authRequest body AuthRequest true "Данные авторизации"
//		@Success		200 {object}	SignInResponse
//		@Failure		400	{object}	errorResponse
//		@Failure		500	{object}	errorResponse
//		@Failure		401	{object}	errorResponse
//		@Router			/auth/ [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.Auth.SignIn"
	log := h.log.With(slog.String("op", op))
	var auth AuthRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "Server error",
			"Server error. Please, try again or later", err.Error())
		return
	}
	err = json.Unmarshal(body, &auth)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest,
			r.Host+r.RequestURI, "Wrong input", "Error parsing body. Please, check your input", err.Error())
		return
	}
	token, err := h.services.SignIn(auth.Username, auth.Password)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) || errors.Is(err, service.ErrUserNotFound) {
			newErrResponse(log, w, http.StatusUnauthorized,
				r.Host+r.RequestURI, "Wrong auth credentials",
				"Incorrect login or password. Please, check your credentials", err.Error())
		} else {
			newErrResponse(log, w, http.StatusInternalServerError,
				r.Host+r.RequestURI, "Server error", "Please, try again or later", err.Error())
		}
		return
	}
	response, err := json.Marshal(SignInResponse{Token: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

type SignUpResponse struct {
	Status int
}

// SignUp godoc
//
//		@Summary		Регистрация
//		@Description	Добавление пользователя
//		@Tags			auth
//		@Accept			json
//		@Produce		json
//	 	@Param			user body domain.User true "Данные регистрации"
//		@Success		200
//		@Failure		400	{object}	errorResponse
//		@Failure		500	{object}	errorResponse
//		@Failure		401	{object}	errorResponse
//		@Router			/signup/ [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.services.SignUp(user)
	if err != nil {
		if errors.Is(err, service.ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp, _ := json.Marshal(SignUpResponse{Status: http.StatusCreated})
	w.Write(resp)
}
