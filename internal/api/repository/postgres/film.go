package postgres

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type FilmPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewFilmPostgres(db *sqlx.DB, log *slog.Logger) *FilmPostgres {
	return &FilmPostgres{db: db, log: log}
}

func (r FilmPostgres) CreateFilm(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) UpdateFilm(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) ListFilms() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) SearchFilm(params repository.SearchFilmParams) ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}
