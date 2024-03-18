package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

func prepare(t *testing.T) (sqlmock.Sqlmock, *sqlx.DB, *ActorPostgres) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	r := NewActorPostgres(dbx, log)

	return mock, dbx, r
}

func TestActorPostgres_CreateActor(t *testing.T) {
	mock, dbx, r := prepare(t)
	defer dbx.Close()

	t.Run("RightCredentials", func(t *testing.T) {
		actor := domain.Actor{
			Id:       1,
			Name:     gofakeit.Name(),
			Birthday: domain.CustomDate(gofakeit.Date()),
			Gender:   1,
		}

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(actor.Id)
		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, actorsTable)).
			WithArgs(actor.Name, actor.Birthday.String(), actor.Gender).WillReturnRows(rows)

		got, err := r.CreateActor(actor)
		assert.NoError(t, err)
		assert.Equal(t, actor.Id, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	/*t.Run("WrongBirthday", func(t *testing.T) {
		actor := domain.Actor{
			Id:       1,
			Name:     gofakeit.Name(),
			Birthday: ("31-02-2001",
			Gender:   1,
		}

		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, actorsTable)).
			WithArgs(actor.Name, actor.Birthday, 2).WillReturnError(errors.New("sql: no rows in result set"))

		got, err := r.CreateActor(actor)
		assert.Error(t, err)
		assert.Equal(t, -1, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})*/

}

func TestActorPostgres_PatchActor(t *testing.T) {
	mock, dbx, r := prepare(t)
	defer dbx.Close()

	t.Run("UpdateAllColumns", func(t *testing.T) {
		actor := domain.Actor{
			Id:       1,
			Name:     gofakeit.Name(),
			Gender:   1,
			Birthday: domain.CustomDate(time.Now()),
		}
		actorInput := domain.ActorInput{
			Id:       1,
			Name:     &actor.Name,
			Gender:   &actor.Gender,
			Birthday: &actor.Birthday,
		}
		rows := sqlmock.NewRows([]string{"id", "name", "birthday", "gender"}).
			AddRow(actor.Id, actor.Name, time.Time(actor.Birthday), actor.Gender)
		mock.ExpectQuery(fmt.Sprintf(`UPDATE %s`, actorsTable)).
			WithArgs(actor.Name, actor.Gender, time.Time(actor.Birthday), actor.Id).WillReturnRows(rows)
		got, err := r.PatchActor(actorInput)
		assert.NoError(t, err)
		assert.Equal(t, actor, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
