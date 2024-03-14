package postgres

import (
	"errors"
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

const (
	uniqueErrCode = "23505"
)

type AuthPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, log *slog.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, log: log}
}

func (r *AuthPostgres) SignUp(user domain.User) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s(username,password_hash) VALUES($1,$2) RETURNING id`, usersTable)
	row := r.db.QueryRowx(query, user.Username, user.PasswordHash)
	var id int
	if err := row.Scan(&id); err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueErrCode {
			return -1, ErrUnique
		}

		return -1, ErrInternal
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1`, usersTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) {
			return user, ErrInternal
		} else {
			return user, ErrNoRows
		}
	}
	return user, nil
}
