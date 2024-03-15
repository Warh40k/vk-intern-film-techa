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

type ActorPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewActorPostgres(db *sqlx.DB, log *slog.Logger) *ActorPostgres {
	return &ActorPostgres{db: db, log: log}
}

func (r ActorPostgres) CreateActor(actor domain.Actor) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s(name, birthday, gender) VALUES($1,$2,$3) RETURNING id`, actorsTable)
	row := r.db.QueryRowx(query, actor.Name, actor.Birthday, actor.Gender)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r ActorPostgres) DeleteActor(id int) error {
	query := fmt.Sprintf(`DELETE FROM %s where id=$1`, actorsTable)
	res, err := r.db.Exec(query, id)
	if err != nil {
		return ErrInternal
	}
	count, err := res.RowsAffected()
	if err != nil {
		return ErrInternal
	}
	if count == 0 {
		return ErrNoRows
	}
	return nil
}

func (r ActorPostgres) UpdateActor(actor domain.Actor) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, gender=$2, birthday=$3 WHERE id=$4`, actorsTable)
	_, err := r.db.Exec(query, actor.Name, actor.Gender, actor.Birthday, actor.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r ActorPostgres) PatchActor(input domain.ActorInput) (domain.Actor, error) {
	var actor domain.Actor

	queryBegin := fmt.Sprintf(`UPDATE %s SET `, actorsTable)
	setVals := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	argId := 1
	if input.Name != nil {
		setVals = append(setVals, "name=$"+strconv.Itoa(argId))
		params = append(params, *input.Name)
		argId++
	}
	if input.Gender != nil {
		setVals = append(setVals, "gender=$"+strconv.Itoa(argId))
		params = append(params, *input.Gender)
		argId++
	}
	if input.Birthday != nil {
		setVals = append(setVals, "birthday=$"+strconv.Itoa(argId))
		params = append(params, time.Time(*input.Birthday))
		argId++
	}

	setString := strings.Join(setVals, ",")
	params = append(params, input.Id)
	query := queryBegin + setString + " WHERE id=$" + strconv.Itoa(argId) +
		" RETURNING *"
	rows := r.db.QueryRowx(query, params...)
	err := rows.StructScan(&actor)

	return actor, err
}

func (r ActorPostgres) ListActors() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}
