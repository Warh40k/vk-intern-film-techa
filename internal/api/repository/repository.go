package repository

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository/postgres"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

//go:generate mockery --all --dry-run=false

type Authorization interface {
	SignUp(user domain.User) (int, error)
	GetUserByUsername(username string) (domain.User, error)
	GetUserById(id int) (domain.User, error)
}

type Actor interface {
	CreateActor(actor domain.Actor) (int, error)
	DeleteActor(id int) error
	UpdateActor(actor domain.Actor) error
	PatchActor(actor domain.ActorInput) (domain.Actor, error)
	ListActors(int) ([]domain.Actor, error)
}

type Film interface {
	CreateFilm(film domain.Film, actorIds []int) (int, error)
	DeleteFilm(id int) error
	UpdateFilm(film domain.Film, actorIds []int) error
	PatchFilm(input domain.NullableFilm, actorIds []int) (domain.Film, error)
	ListFilms(sortBy, sortDir string) ([]domain.Film, error)
	SearchFilm(query string) ([]domain.Film, error)
	ListFilmsByActor(sortBy, sortDir string, actorId int) ([]domain.Film, error)
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
