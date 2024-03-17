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

func (r FilmPostgres) PatchFilm(film domain.FilmInput) (domain.Film, error) {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) CreateFilm(film domain.Film, actorIds []int) (int, error) {
	var filmId int
	const method = "Films.Repository.CreateFilm"
	log := r.log.With(slog.String("method", method))

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, ErrInternal
	}

	createFilmQuery := fmt.Sprintf(`INSERT INTO %s(title, description, released, rating) 
		VALUES($1,$2,$3,$4) RETURNING id`, filmsTable)
	row := tx.QueryRowx(createFilmQuery, film.Title, film.Description, film.Released.String(), film.Rating)
	if err = row.Scan(&filmId); err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return 0, ErrInternal
	}

	addActorsStmt, err := tx.Preparex(
		fmt.Sprintf(`INSERT INTO %s(film_id, actor_id) VALUES($1,$2)`, filmsActorsTable))
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return 0, ErrInternal
	}

	for _, actorId := range actorIds {
		_, err = addActorsStmt.Exec(filmId, actorId)
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())
			return 0, ErrInternal
		}
	}

	err = tx.Commit()

	return filmId, err
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
