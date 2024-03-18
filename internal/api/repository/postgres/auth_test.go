package postgres

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"regexp"
	"testing"
)

func TestAuthPostgres_SignUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	r := NewAuthPostgres(dbx, log)

	type mockBehaviour func(input domain.User, id int)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		input         domain.User
		want          int
		wantErr       bool
	}{
		{
			name: "ValidCredentials",
			mockBehaviour: func(input domain.User, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, usersTable)).
					WithArgs(input.Username, input.PasswordHash).WillReturnRows(rows)
			},
			input: domain.User{
				Username:     "test",
				PasswordHash: "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO",
			},
			want: 1,
		},
		{
			name: "UniqueViolation",
			mockBehaviour: func(input domain.User, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, usersTable)).
					WithArgs(input.Username, input.PasswordHash).WillReturnRows(rows)
				mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s`, usersTable)).
					WithArgs(input.Username, input.PasswordHash).WillReturnError(pgx.PgError{
					Code: uniqueErrCode,
				})
			},
			input: domain.User{
				Username:     "test",
				PasswordHash: "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO",
			},
			want: 1,
		},
	}

	t.Run(tests[0].name, func(t *testing.T) {
		tests[0].mockBehaviour(tests[0].input, tests[0].want)
		got, err := r.SignUp(tests[0].input)
		assert.NoError(t, err)
		assert.Equal(t, tests[0].want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run(tests[1].name, func(t *testing.T) {
		tests[1].mockBehaviour(tests[1].input, tests[1].want)
		got, err := r.SignUp(tests[1].input)
		assert.NoError(t, err)
		assert.Equal(t, tests[1].want, got)

		got, err = r.SignUp(tests[1].input)

		assert.ErrorIs(t, err, ErrUnique)
		assert.Equal(t, -1, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthPostgres_GetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	r := NewAuthPostgres(dbx, log)

	user := domain.User{
		Id:           1,
		Username:     gofakeit.Username(),
		PasswordHash: "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO",
		Role:         1,
	}

	t.Run("RightCredentials", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role"}).
			AddRow(user.Id, user.Username, user.PasswordHash, user.Role)
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT * FROM %s`), usersTable)).
			WithArgs(user.Username).WillReturnRows(rows)

		got, err := r.GetUserByUsername(user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UserNotExist", func(t *testing.T) {
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT * FROM %s`), usersTable)).
			WithArgs(user.Username).WillReturnError(errors.New("sql: no rows in result set"))

		_, err := r.GetUserByUsername(user.Username)
		assert.ErrorIs(t, ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestAuthPostgres_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	r := NewAuthPostgres(dbx, log)

	user := domain.User{
		Id:           1,
		Username:     gofakeit.Username(),
		PasswordHash: "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO",
		Role:         1,
	}

	t.Run("RightCredentials", func(t *testing.T) {

		gofakeit.Username()
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role"}).
			AddRow(user.Id, user.Username, user.PasswordHash, user.Role)
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT * FROM %s`), usersTable)).
			WithArgs(user.Id).WillReturnRows(rows)

		got, err := r.GetUserById(user.Id)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UserNotExist", func(t *testing.T) {
		gofakeit.Username()
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta(`SELECT * FROM %s`), usersTable)).
			WithArgs(user.Id).WillReturnError(errors.New("sql: no rows in result set"))

		_, err := r.GetUserById(user.Id)
		assert.ErrorIs(t, ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
