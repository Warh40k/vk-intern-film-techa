package domain

type User struct {
	Id           int    `json:"-" db:"id"`
	Username     string `json:"username" db:"username" validate:"required"`
	Password     string `json:"password" db:"-" validate:"required"`
	PasswordHash string `json:"-" db:"password_hash"`
	Role         int8   `json:"-" db:"role"`
}
