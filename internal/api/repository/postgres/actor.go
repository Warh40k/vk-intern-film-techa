package postgres

import (
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
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
	row := r.db.QueryRowx(query, actor.Name, actor.Birthday.Date, actor.Gender)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r ActorPostgres) DeleteActor(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r ActorPostgres) UpdateActor(actor domain.Actor) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, gender=$2, birthday=$3 WHERE id=$4`, actorsTable)
	_, err := r.db.Exec(query, actor.Name, actor.Gender, actor.Birthday.Date, actor.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r ActorPostgres) PatchActor(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (r ActorPostgres) ListActors() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}
