package postgres

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"os"
	"regexp"
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

func TestActorPostgres_UpdateActor(t *testing.T) {
	mock, dbx, r := prepare(t)
	defer dbx.Close()

	t.Run("RightCredentials", func(t *testing.T) {
		actor := domain.Actor{
			Id:       1,
			Name:     gofakeit.Name(),
			Gender:   1,
			Birthday: domain.CustomDate(time.Now()),
		}
		mock.ExpectExec(fmt.Sprintf(`UPDATE %s`, actorsTable)).
			WithArgs(actor.Name, actor.Gender, actor.Birthday.String(), actor.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := r.UpdateActor(actor)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("NoSuchId", func(t *testing.T) {
		actor := domain.Actor{
			Id:       1,
			Name:     gofakeit.Name(),
			Gender:   1,
			Birthday: domain.CustomDate(time.Now()),
		}
		mock.ExpectExec(fmt.Sprintf(`UPDATE %s`, actorsTable)).
			WithArgs(actor.Name, actor.Gender, actor.Birthday.String(), actor.Id).
			WillReturnResult(sqlmock.NewResult(1, 0))
		err := r.UpdateActor(actor)
		assert.ErrorIs(t, err, ErrNoRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestActorPostgres_ListActors(t *testing.T) {
	mock, dbx, r := prepare(t)
	defer dbx.Close()

	t.Run("GetAll", func(t *testing.T) {
		actors := []domain.Actor{
			{
				Id:       1,
				Name:     gofakeit.Name(),
				Gender:   1,
				Birthday: domain.CustomDate(time.Now()),
				/*Films: []domain.Film{
					{
						Id:          3,
						Title:       "test",
						Description: "test",
						Released:    domain.CustomDate(gofakeit.Date()),
						Rating:      10,
					},
					{
						Id:          filmId,
						Title:       "test2",
						Description: "test2",
						Released:    domain.CustomDate(gofakeit.Date()),
						Rating:      7,
					},
				},*/
			},
		}
		rows := sqlmock.NewRows([]string{"id", "name", "birthday", "gender"}).
			AddRows([][]driver.Value{{actors[0].Id, actors[0].Name, time.Time(actors[0].Birthday), actors[0].Gender}}...)
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT a.* FROM %s a`), actorsTable)).
			WithoutArgs().WillReturnRows(rows)
		got, err := r.ListActors(-1)
		assert.NoError(t, err)
		assert.Equal(t, got, actors)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetByFilmId", func(t *testing.T) {
		filmId := rand.Int()
		actors := []domain.Actor{
			{
				Id:       1,
				Name:     gofakeit.Name(),
				Gender:   1,
				Birthday: domain.CustomDate(time.Now()),
			},
		}
		rows := sqlmock.NewRows([]string{"id", "name", "birthday", "gender"}).
			AddRows([][]driver.Value{{actors[0].Id, actors[0].Name, time.Time(actors[0].Birthday), actors[0].Gender}}...)
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT a.* FROM %s a`), actorsTable)).
			WithArgs(filmId).WillReturnRows(rows)
		got, err := r.ListActors(filmId)
		assert.NoError(t, err)
		assert.Equal(t, got, actors)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
