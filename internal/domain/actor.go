package domain

type Actor struct {
	Id       int        `json:"id" db:"id"`
	Name     string     `json:"name" db:"name" validate:"required"`
	Gender   int        `json:"gender" db:"gender" validate:"required,oneof=0 1 2 9"` // ISO/IEC 5218
	Birthday CustomDate `json:"birthday" db:"birthday" validate:"required"`
	Films    []Film     `json:"films" db:"-"`
}

type ActorInput struct {
	Id       int         `json:"-"`
	Name     *string     `json:"name" validate:"omitempty,gt=0"`
	Gender   *int        `json:"gender" validate:"omitempty,oneof=0 1 2 9"`
	Birthday *CustomDate `json:"birthday" validate:"omitempty"`
}
