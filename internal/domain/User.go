package domain

import "time"

type User struct {
	Id           int       `json:"id,omitempty" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required"`
	Username     string    `json:"username" db:"username" validate:"required"`
	Password     string    `json:"password" db:"-" validate:"required"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Created      time.Time `json:"created" db:"created"`
}
