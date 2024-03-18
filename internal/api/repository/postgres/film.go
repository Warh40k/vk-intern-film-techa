package postgres

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type FilmPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func (r FilmPostgres) ListFilms(sortBy, sortDir string) ([]domain.Film, error) {
	var films []domain.Film
	query := fmt.Sprintf(`SELECT * from %s ft ORDER BY %s %s`, filmsTable, sortBy, sortDir)
	err := r.db.Select(&films, query)

	return films, err
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
			return ErrUnique
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
		params = append(params, time.Time(*input.Released))
		argId++
	}
	if input.Rating != nil {
		setVals = append(setVals, "rating=$"+strconv.Itoa(argId))
		params = append(params, *input.Rating)
		argId++
	}

	tx, err := r.db.Beginx()
	if err != nil {
		log.Error(err.Error())
		return film, err
	}

	if argId != 1 {
		setString := strings.Join(setVals, ",")
		params = append(params, input.Id)
		query := queryBegin + setString + " WHERE id=$" + strconv.Itoa(argId) +
			" RETURNING *"
		rows := tx.QueryRowx(query, params...)
		err = rows.StructScan(&film)
		if err != nil {
			log.Error(err.Error())
			tx.Rollback()
			return film, ErrNoRows
		}
	}

	err = r.updateActorsList(tx, input.Id, actorIds)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		return film, err
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
	err = r.updateActorsList(tx, filmId, actorIds)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return 0, err
	}

	return filmId, tx.Commit()
}

func (r FilmPostgres) DeleteFilm(id int) error {
	const method = "Films.Repository.DeleteFilm"
	log := r.log.With(slog.String("method", method))

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", filmsTable)
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error(err.Error())
		return ErrInternal
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return ErrInternal
	}
	if count == 0 {
		log.Info("No rows affected")
		return ErrNoRows
	}

	return nil
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
		return err
	}

	return tx.Commit()
}

func (r FilmPostgres) ListFilmsByActor(sortBy, sortDir string, actorId int) ([]domain.Film, error) {
	var films []domain.Film
	query := fmt.Sprintf(`SELECT ft.* from %s ft INNER JOIN %s fa ON ft.id = fa.film_id 
                                     WHERE fa.actor_id = $1 ORDER BY %s %s`,
		filmsTable, filmsActorsTable, sortBy, sortDir)

	err := r.db.Select(&films, query, actorId)

	return films, err
}

func (r FilmPostgres) SearchFilm(searchQuery string) ([]domain.Film, error) {
	var films []domain.Film
	query := fmt.Sprintf(`SELECT f.* FROM %s f 
           INNER JOIN %s fa ON f.id = fa.film_id 
           INNER JOIN %s a ON a.id = fa.actor_id 
           WHERE f.title LIKE $1 OR a.name LIKE $1 GROUP BY f.id`, filmsTable, filmsActorsTable, actorsTable)
	like := fmt.Sprintf("%%%s%%", searchQuery)
	err := r.db.Select(&films, query, like)
	return films, err
}
