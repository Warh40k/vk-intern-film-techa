package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func prepareFilmTest(t *testing.T) (sqlmock.Sqlmock, *sqlx.DB, *FilmPostgres) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	r := NewFilmPostgres(dbx, log)

	return mock, dbx, r
}

func TestFilmPostgres_CreateFilm(t *testing.T) {
	mock, dbx, r := prepareFilmTest(t)
	defer dbx.Close()

	t.Run("RightCredentials", func(t *testing.T) {
		film := domain.Film{
			Id:       1,
			Title:    gofakeit.JobTitle(),
			Released: domain.CustomDate(gofakeit.Date()),
			Rating:   5,
		}

		actorIds := []int{1, 2, 3}

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(film.Id)
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, filmsTable)).
			WithArgs(film.Title, film.Description, film.Released.String(), film.Rating).WillReturnRows(rows)
		mock.ExpectPrepare(fmt.Sprintf("INSERT INTO %s", filmsActorsTable))
		for _, actorId := range actorIds {
			mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", filmsActorsTable)).
				WithArgs(film.Id, actorId).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()
		got, err := r.CreateFilm(film, actorIds)
		assert.NoError(t, err)
		assert.Equal(t, film.Id, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UniqueActorIds", func(t *testing.T) {
		film := domain.Film{
			Id:       1,
			Title:    gofakeit.JobTitle(),
			Released: domain.CustomDate(gofakeit.Date()),
			Rating:   5,
		}

		actorIds := []int{2, 2}

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(film.Id)
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, filmsTable)).
			WithArgs(film.Title, film.Description, film.Released.String(), film.Rating).WillReturnRows(rows)
		mock.ExpectPrepare(fmt.Sprintf("INSERT INTO %s", filmsActorsTable))

		mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", filmsActorsTable)).
			WithArgs(film.Id, actorIds[0]).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", filmsActorsTable)).
			WithArgs(film.Id, actorIds[1]).
			WillReturnError(pgx.PgError{Code: uniqueErrCode})
		mock.ExpectRollback()
		_, err := r.CreateFilm(film, actorIds)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUnique)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
