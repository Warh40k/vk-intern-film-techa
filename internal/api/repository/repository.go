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
	CreateActor(actor domain.Actor) error
	DeleteActor(id int) error
	UpdateActor(actor domain.Actor) error
	ListActors() ([]domain.Actor, error)
}

type SearchFilmParams struct {
	Filmname  string
	Actorname string
}

type Film interface {
	CreateFilm(actor domain.Actor) error
	DeleteFilm(id int) error
	UpdateFilm(actor domain.Actor) error
	ListFilms() ([]domain.Actor, error)
	SearchFilm(params SearchFilmParams) ([]domain.Film, error)
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
