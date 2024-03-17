package repository

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository/postgres"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Authorization interface {
	SignUp(user domain.User) (int, error)
	GetUserByUsername(username string) (domain.User, error)
}

type Actor interface {
	CreateActor(actor domain.Actor) (int, error)
	DeleteActor(id int) error
	UpdateActor(actor domain.Actor) error
	PatchActor(actor domain.ActorInput) (domain.Actor, error)
	ListActors() ([]domain.Actor, error)
}

type Film interface {
	CreateFilm(film domain.Film) (int, error)
	DeleteFilm(id int) error
	UpdateFilm(film domain.Film) error
	PatchFilm(film domain.FilmInput) error
	ListFilms(sortBy, sortDir string) ([]domain.Film, error)
	SearchFilm(film string, actor string) ([]domain.Film, error)
}

type Repository struct {
	Authorization
	Actor
	Film
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db, log),
		Film:          postgres.NewFilmPostgres(db, log),
		Actor:         postgres.NewActorPostgres(db, log),
	}
}
