package domain

type Actor struct {
	Id       int        `json:"id" db:"id"`
	Name     string     `json:"name" db:"name" validate:"required"`
	Gender   int        `json:"gender" db:"gender" validate:"required,oneof=0 1 2 9"` // ISO/IEC 5218
	Birthday CustomDate `json:"birthday" db:"birthday" validate:"required"`
}
