package postgres

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strconv"
	"strings"
)

type FilmPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewFilmPostgres(db *sqlx.DB, log *slog.Logger) *FilmPostgres {
	return &FilmPostgres{db: db, log: log}
}

func (r FilmPostgres) updateActorsList(tx *sqlx.Tx, filmId int, actorIds []int) error {
	const method = "Films.Repository.updateActorsList"
	log := r.log.With(slog.String("method", method))

	clearOldActors := fmt.Sprintf(`DELETE FROM %s where film_id=$1`, filmsActorsTable)
	_, err := tx.Exec(clearOldActors, filmId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	addActorsStmt, err := tx.Preparex(
		fmt.Sprintf(`INSERT INTO %s(film_id, actor_id) VALUES($1,$2)`, filmsActorsTable))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	for _, actorId := range actorIds {
		_, err = addActorsStmt.Exec(filmId, actorId)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (r FilmPostgres) PatchFilm(input domain.NullableFilm, actorIds []int) (domain.Film, error) {
	const method = "Films.Repository.PatchFilm"
	log := r.log.With(slog.String("method", method))

	var film domain.Film

	queryBegin := fmt.Sprintf(`UPDATE %s SET `, filmsTable)
	setVals := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	argId := 1
	if input.Title != nil {
		setVals = append(setVals, "title=$"+strconv.Itoa(argId))
		params = append(params, *input.Title)
		argId++
	}
	if input.Description != nil {
		setVals = append(setVals, "description=$"+strconv.Itoa(argId))
		params = append(params, *input.Description)
		argId++
	}
	if input.Released != nil {
		setVals = append(setVals, "released=$"+strconv.Itoa(argId))
		params = append(params, input.Released.String())
		argId++
	}
	if input.Rating != nil {
		setVals = append(setVals, "rating=$"+strconv.Itoa(argId))
		params = append(params, *input.Rating)
		argId++
	}

	setString := strings.Join(setVals, ",")
	params = append(params, input.Id)
	query := queryBegin + setString + " WHERE id=$" + strconv.Itoa(argId) +
		" RETURNING *"
	tx, err := r.db.Beginx()
	if err != nil {
		log.Error(err.Error())
		return film, err
	}
	rows := tx.QueryRowx(query, params...)
	err = rows.StructScan(&film)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return film, ErrNoRows
	}

	err = r.updateActorsList(tx, input.Id, actorIds)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return film, ErrInternal
	}

	return film, tx.Commit()
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

	err = r.updateActorsList(tx, film.Id, actorIds)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return ErrInternal
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
