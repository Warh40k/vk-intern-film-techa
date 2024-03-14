package handler

import (
	"encoding/json"
	"errors"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required,gte=1,lte=128"`
	Password string `json:"password" validate:"required,gte=8,lte=128"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var auth AuthRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &auth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.services.SignIn(auth.Username, auth.Password)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			w.WriteHeader(http.StatusUnauthorized)
		} else if errors.Is(err, service.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
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
}
