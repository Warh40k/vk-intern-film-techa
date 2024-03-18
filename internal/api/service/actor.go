package service

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"log/slog"
)

type ActorService struct {
	repos repository.Actor
	log   *slog.Logger
}

func (s *ActorService) PatchActor(actor domain.ActorInput) (domain.Actor, error) {
	return s.repos.PatchActor(actor)
}

func NewActorService(repos repository.Actor, log *slog.Logger) *ActorService {
	return &ActorService{repos: repos, log: log}
}

func (s *ActorService) CreateActor(actor domain.Actor) (int, error) {
	return s.repos.CreateActor(actor)
}

func (s *ActorService) DeleteActor(id int) error {
	return s.repos.DeleteActor(id)
}

func (s *ActorService) UpdateActor(actor domain.Actor) error {
	return s.repos.UpdateActor(actor)
}

func (s *ActorService) ListActors(filmId int) (actors []domain.Actor, err error) {
	return s.repos.ListActors(filmId)
}
