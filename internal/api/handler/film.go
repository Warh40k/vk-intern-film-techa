package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
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

	films, err := h.services.ListFilms(sortParams[0], sortParams[1], -1)
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

type FilmInput struct {
	domain.Film `json:"film"`
	ActorIds    []int `json:"actorIds,omitempty"`
}

func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.CreateFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	var input FilmInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "data parse error",
			"Failed to parse data. Please, check your input", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "validation error",
			"Couldn't validate input fields. Please, fix input and try again", vErr.Error())
		return
	}

	input.Id, err = h.services.CreateFilm(input.Film, input.ActorIds)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "save film error",
			"Failed to save film. Please, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(input)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete method"))

}

type PatchFilmInput struct {
	domain.NullableFilm `json:"film"`
	ActorIds            []int `json:"actorIds"`
}

func (h *Handler) PatchFilm(w http.ResponseWriter, r *http.Request) {
	// TODO: Дропать актеров, и добавлять новое значение (если оно есть)
	const method = "Handlers.Film.PatchFilm"
	log := h.log.With(
		slog.String("method", method),
	)
	var input PatchFilmInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "data parse error",
			"Failed to parse data. Please, check your input", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "validation error",
			"Couldn't validate input fields. Please, fix input and try again", vErr.Error())
		return
	}

	input.Id, err = strconv.Atoi(r.PathValue("film_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "param error",
			"Fsailed to get film id. Please, check your input and try again", err.Error())
		return
	}

	film, err := h.services.PatchFilm(input.NullableFilm, input.ActorIds)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "save error",
			"failed to save data, try again later", err.Error())
		return
	}

	resp, err := json.Marshal(map[string]interface{}{"film": film, "actorIds": input.ActorIds})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	// TODO: разобраться с вводом несуществующих Id актеров
	const method = "Handlers.Film.UpdateFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	var input FilmInput
	var err error
	input.Film.Id, err = strconv.Atoi(r.PathValue("film_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "param error",
			"failed to get user id. Please, check your input and try again", err.Error())
		return
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "json parse error",
			"Failed to parse json. Please, check your input", err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "validation error",
			"Couldn't validate input fields. Please, check input and try again", vErr.Error())
		return
	}

	err = h.services.UpdateFilm(input.Film, input.ActorIds)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "save error",
			"Failed to save data, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(input.Film)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
