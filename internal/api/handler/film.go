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

// ListFilms godoc
//
//		@Summary		Список фильмов
//		@Description	Получить список фильмов
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			sortby query string true "Поле и направление сортировки" example(rating.desc)
//		@Success		200	{array}		domain.Film
//		@Failure		400	{object}	errorResponse
//		@Failure		500	{object}	errorResponse
//		@Router			/films/ [get]
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

	for i := range films {
		films[i].Actors, err = h.services.ListActors(films[i].Id)
		if err != nil {
			newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "server error",
				"Failed to get data. Please, try again later", err.Error())
			return
		}
		if films[i].Actors == nil {
			films[i].Actors = []domain.Actor{}
		}
	}

	if films == nil {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal([]domain.Film{})
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(films)
	w.Write(resp)
}

// SearchFilm godoc
//
//		@Summary		Поиск фильмов
//		@Description	Поиск фильмов по фрагменту названия фильма или имени актера
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			query query string true "Поисковый запрос" example("Avatar")
//		@Success		200	{array}		domain.Film
//		@Failure		400	{object}	errorResponse
//		@Failure		404	{object}	errorResponse
//		@Failure		500	{object}	errorResponse
//		@Router			/films/search [get]
func (h *Handler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.SearchFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	query := r.URL.Query().Get("query")
	if query == "" {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "input error",
			"Search query is empty", "Search query is empty")
		return
	}
	films, err := h.services.SearchFilm(query)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "server error",
			"Failed to get data. Please, try again later", err.Error())
		return
	}

	for i := range films {
		films[i].Actors, err = h.services.ListActors(films[i].Id)
		if err != nil {
			newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "server error",
				"Failed to get data. Please, try again later", err.Error())
			return
		}
		if films[i].Actors == nil {
			films[i].Actors = []domain.Actor{}
		}
	}

	if films == nil {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal([]domain.Film{})
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(films)
	w.Write(resp)
}

type filmInput struct {
	domain.Film `json:"film"`
	ActorIds    []int `json:"actorIds,omitempty"`
}

// CreateFilm godoc
//
//		@Summary		Добавить фильм
//		@Description	Добавить информацию по фильму
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			input body filmInput true "Информация о фильму" example("Avatar")
//		@Success		200	{object}	domain.Film
//		@Failure		400	{object}	errorResponse
//		@Failure		404	{object}	errorResponse
//		@Failure		500	{object}	errorResponse
//		@Router			/films/ [post]
func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.CreateFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	var input filmInput
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

// DeleteFilm godoc
//
//		@Summary		Удалить фильм
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			film_id path int true "ИД фильма" example(10)
//		@Success		200
//		@Failure		400	{object}	errorResponse
//		@Router			/films/{film_id}/ [delete]
func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Film.PatchFilm"
	log := h.log.With(
		slog.String("method", method),
	)
	filmId, err := strconv.Atoi(r.PathValue("film_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "input error",
			"Incorrect film id. Please, check your input", err.Error())
		return
	}

	err = h.services.DeleteFilm(filmId)
	if err != nil {
		// TODO: улучшить обработку ошибок
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "server error",
			"Internal error. Please, try again later", err.Error())
		return
	}
}

type PatchFilmInput struct {
	domain.NullableFilm `json:"film"`
	ActorIds            []int `json:"actorIds"`
}

// PatchFilm godoc
//
//		@Summary		Редактировать фильм
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			filmInput body PatchFilmInput true "Данные для обновления"
//	 	@Param			film_id path int true "ИД фильма"
//		@Success		200 {object}	domain.Film
//		@Failure		400	{object}	errorResponse
//		@Router			/films/{film_id}/ [patch]
func (h *Handler) PatchFilm(w http.ResponseWriter, r *http.Request) {
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

// UpdateFilm godoc
//
//		@Summary		Обновить фильм
//		@Description	Полная замена фильма
//		@Tags			films
//		@Accept			json
//		@Produce		json
//	 	@Param			film body domain.Film true "Данные фильма"
//	 	@Param			film_id path int true "ИД фильма"
//		@Success		200 {object}	domain.Film
//		@Failure		400	{object}	errorResponse
//		@Router			/films/{film_id}/ [put]
func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	// TODO: разобраться с вводом несуществующих Id актеров
	const method = "Handlers.Film.UpdateFilm"
	log := h.log.With(
		slog.String("method", method),
	)

	var input filmInput
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
