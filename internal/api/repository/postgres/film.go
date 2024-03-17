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

func (r FilmPostgres) PatchFilm(film domain.PatchFilmInput, actorIds []int) (domain.Film, error) {
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

	return filmId, tx.Commit()
}

func (r FilmPostgres) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r FilmPostgres) UpdateFilm(film domain.Film, actorIds []int) error {
	const method = "Films.Repository.UpdateFilm"
	log := r.log.With(slog.String("method", method))

	tx, err := r.db.Beginx()
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return ErrInternal
	}
	modifyFilmInfo := fmt.Sprintf(`UPDATE %s SET title=$1, description=$2, released=$3, rating=$4 
          WHERE id=$5`, filmsTable)
	_, err = tx.Exec(modifyFilmInfo, film.Title, film.Description, film.Released.String(), film.Rating, film.Id)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return err
	}

	clearOldActors := fmt.Sprintf(`DELETE FROM %s where film_id=$1`, filmsActorsTable)
	_, err = tx.Exec(clearOldActors, film.Id)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return ErrInternal
	}

	insertActors, err := tx.Preparex(
		fmt.Sprintf(`INSERT INTO %s(film_id, actor_id) VALUES($1, $2)`, filmsActorsTable))
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return ErrInternal
	}

	for _, actorId := range actorIds {
		_, err = insertActors.Exec(film.Id, actorId)
		if err != nil {
			log.Error(err.Error())
			tx.Rollback()
			return ErrInternal
		}
	}

	return tx.Commit()
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
