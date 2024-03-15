package domain

type User struct {
	Id           int    `json:"id,omitempty" db:"id"`
	Username     string `json:"username" db:"username" validate:"required"`
	Password     string `json:"password" db:"-" validate:"required"`
	PasswordHash string `json:"-" db:"password_hash"`
}