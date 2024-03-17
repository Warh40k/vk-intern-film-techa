package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
)

func validateSortParams(sortParams []string) ([]string, error) {
	if sortParams[0] == "" {
		sortParams[0] = sortRating
	}
	if len(sortParams) == 1 {
		sortParams = append(sortParams, descSort)
	}

	if sortParams[0] != sortTitle && sortParams[0] != sortRating && sortParams[0] != sortReleased {
		return sortParams, fmt.Errorf("unsupported sorting column %q", sortParams[0])
	}
	if sortParams[1] != ascSort && sortParams[1] != descSort {
		return sortParams, fmt.Errorf("unsupported sorting direction %q", sortParams[1])
	}

	return sortParams, nil
}

func (h *Handler) ListFilms(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.ListFilms"
	log := h.log.With(
		slog.String("method", method),
	)

	sortParams := strings.Split(r.URL.Query().Get("sortby"), ".")
	sortParams, err := validateSortParams(sortParams)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "sort error", err.Error(), err.Error())
		return
	}

	films, err := h.services.ListFilms(sortParams[0], sortParams[1])
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "sort error",
			"Failed to get films. Please, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(films)
	w.Write(resp)
}

func (h *Handler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get method"))
}

func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.CreateFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	var film domain.Film
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "data parse error",
			"Failed to parse data. Please, check your input", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(film)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "validation error",
			"Couldn't validate input fields. Please, fix input and try again", vErr.Error())
		return
	}

	film.Id, err = h.services.CreateFilm(film)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "save film error",
			"Failed to save film. Please, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(film)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete method"))

}

func (h *Handler) PatchFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Patch method"))

}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Put method"))

}
