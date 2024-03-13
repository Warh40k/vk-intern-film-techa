package postgres

import (
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

func (r ActorPostgres) CreateActor(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (r ActorPostgres) DeleteActor(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r ActorPostgres) UpdateActor(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (r ActorPostgres) ListActors() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}
