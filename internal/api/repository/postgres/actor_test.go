package postgres

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestActorPostgres_CreateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	r := NewActorPostgres(dbx, log)

	actor := domain.Actor{
		Id:       1,
		Name:     gofakeit.Name(),
		Birthday: domain.CustomDate(gofakeit.Date()),
		Gender:   1,
	}

	t.Run("RightCredentials", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(actor.Id)
		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, actorsTable)).
			WithArgs(actor.Name, actor.Birthday.String(), actor.Gender).WillReturnRows(rows)

		got, err := r.CreateActor(actor)
		assert.NoError(t, err)
		assert.Equal(t, actor.Id, got)
	})

	t.Run("WrongBirthday", func(t *testing.T) {
		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, actorsTable)).
			WithArgs(actor.Name, "31-02-2001").WillReturnError(errors.New("sql: no rows in result set"))

		got, err := r.CreateActor(actor)
		assert.Error(t, err)
		assert.Equal(t, -1, got)
	})

}
