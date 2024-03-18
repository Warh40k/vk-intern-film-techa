package handler

import (
	"encoding/json"
	"errors"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) ListActors(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Actor.ListActors"
	log := h.log.With(
		slog.String("method", method),
	)

	actors, err := h.services.ListActors(-1)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "server error",
			"Failed to get actors list. Please, try again later", err.Error())
		return
	}

	for i := range actors {
		actors[i].Films, err = h.services.ListFilms(sortRating, descSort, actors[i].Id)
		if err != nil {
			newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "server error",
				"Failed to get actors list. Please, try again later", err.Error())
			return
		}
		if actors[i].Films == nil {
			actors[i].Films = []domain.Film{}
		}
	}

	if actors == nil {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal([]domain.Actor{})
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(actors)
	w.Write(resp)
}

func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Actor.CreateActor"
	log := h.log.With(
		slog.String("method", method),
	)
	var actor domain.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "json parse error",
			"Failed to parse json. Please, check your input", err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(actor)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "Validation error",
			"Couldn't validate input fields. Please, fix input and try again", vErr.Error())
		return
	}

	actor.Id, err = h.services.CreateActor(actor)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "Create actor error",
			"Failed to create new actor. Please, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(actor)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Actor.DeleteActor"
	log := h.log.With(
		slog.String("method", method),
	)
	id, err := strconv.Atoi(r.PathValue("actor_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "param error",
			"Failed to get actor id. Please, check your input", err.Error())
		return
	}
	log = h.log.With(slog.Int("actor id", id))

	err = h.services.DeleteActor(id)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "delete error",
			"Specified actor not found", err.Error())
		return
	}
}

func (h *Handler) PatchActor(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Actor.PatchActor"
	log := h.log.With(
		slog.String("method", method),
	)
	var input domain.ActorInput
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

	input.Id, err = strconv.Atoi(r.PathValue("actor_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "error getting actor id",
			"failed to get user id. Please, check your input and try again", err.Error())
		return
	}

	actor, err := h.services.PatchActor(input)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "patch error",
			"failed to save data, try again later", err.Error())
		return
	}

	resp, err := json.Marshal(actor)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	const method = "Handlers.Actor.UpdateActor"
	log := h.log.With(
		slog.String("method", method),
	)

	var actor domain.Actor
	var err error
	actor.Id, err = strconv.Atoi(r.PathValue("actor_id"))
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "error getting actor id",
			"failed to get user id. Please, check your input and try again", err.Error())
		return
	}
	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "json parse error",
			"Failed to parse json. Please, check your input", err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(actor)
	if err != nil {
		var vErr validator.ValidationErrors
		errors.As(err, &vErr)
		newErrResponse(log, w, http.StatusBadRequest, r.Host+r.RequestURI, "validation error",
			"Couldn't validate input fields. Please, fix input and try again", vErr.Error())
		return
	}

	err = h.services.UpdateActor(actor)
	if err != nil {
		newErrResponse(log, w, http.StatusInternalServerError, r.Host+r.RequestURI, "error updating actor",
			"Failed to save data, try again later", err.Error())
		return
	}

	resp, _ := json.Marshal(actor)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
