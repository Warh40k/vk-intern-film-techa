package service

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"log/slog"
)

type ActorService struct {
	repos repository.Film
	log   *slog.Logger
}

func NewActorService(repos repository.Film, log *slog.Logger) *ActorService {
	return &ActorService{repos: repos, log: log}
}

func (s *ActorService) CreateActor(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (s *ActorService) DeleteActor(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s *ActorService) UpdateActor(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (s *ActorService) ListActors() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}
