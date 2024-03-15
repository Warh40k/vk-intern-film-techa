package service

import (
	"errors"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"log/slog"
)

//go:generate mockery --all --dry-run=false

var (
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal errors")
	ErrUnauthorized = errors.New("not authorized")
)

type Service struct {
	Authorization
	Actor
	Film
}

type Authorization interface {
	SignUp(user domain.User) error
	SignIn(username, password string) (string, error)
}

type Actor interface {
	CreateActor(actor domain.Actor) (int, error)
	DeleteActor(id int) error
	UpdateActor(actor domain.Actor) error
	PatchActor(actor domain.ActorInput) (domain.Actor, error)
	ListActors() ([]domain.Actor, error)
}

type Film interface {
	CreateFilm(actor domain.Film) error
	DeleteFilm(id int) error
	UpdateFilm(actor domain.Film) error
	PatchFilm(actor domain.Film) error
	ListFilms() ([]domain.Film, error)
	SearchFilm(film string, actor string) ([]domain.Film, error)
}

func NewService(repos *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos, log),
		Actor:         NewActorService(repos, log),
		Film:          NewFilmService(repos, log),
	}
}
