package postgres

import (
	"fmt"
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

func (r FilmPostgres) PatchFilm(film domain.FilmInput) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) CreateFilm(film domain.Film) (int, error) {
	var id int
	const method = "Films.Repository.CreateFilm"
	log := r.log.With(slog.String("method", method))

	query := fmt.Sprintf(`INSERT INTO %s(title, description, released, rating) 
		VALUES($1,$2,$3,$4) RETURNING id`, filmsTable)
	row := r.db.QueryRowx(query, film.Title, film.Description, film.Released.String(), film.Rating)
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
	query := fmt.Sprintf(`UPDATE %s SET title=$1, description=$2, released=$3, rating=$4 
          WHERE id=$5`, filmsTable)
	_, err := r.db.Exec(query, film.Title, film.Description, film.Released.String(), film.Rating, film.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r FilmPostgres) ListFilms(sortBy, sortDir string) ([]domain.Film, error) {
	var films []domain.Film
	query := fmt.Sprintf(`SELECT * from %s ft ORDER BY %s %s`, filmsTable, sortBy, sortDir)
	err := r.db.Select(&films, query)

	return films, err
}

func (r FilmPostgres) SearchFilm(query string) ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}
