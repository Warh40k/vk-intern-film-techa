package postgres

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type FilmPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func (r FilmPostgres) PatchFilm(film domain.FilmInput) error {
	//TODO implement me
	panic("implement me")
}

func NewFilmPostgres(db *sqlx.DB, log *slog.Logger) *FilmPostgres {
	return &FilmPostgres{db: db, log: log}
}

func (r FilmPostgres) CreateFilm(film domain.Film) (int, error) {
	var id int
	const method = "Films.Repository.CreateFilm"
	log := r.log.With(slog.String("method", method))

	query := fmt.Sprintf(`INSERT INTO %s(title, description, released, rating) 
		VALUES($1,$2,$3,$4) RETURNING id`, filmsTable)
	row := r.db.QueryRowx(query, film.Title, film.Description, time.Time(film.Released), film.Rating)
	if err := row.Scan(&id); err != nil {
		log.Error(err.Error())
		return id, ErrInternal
	}
	return id, nil
}

func (r FilmPostgres) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) UpdateFilm(film domain.Film) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) ListFilms() ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) SearchFilm(film string, actor string) ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}
