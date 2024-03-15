package domain

type Film struct {
	Id          int        `json:"id,omitempty" db:"id"`
	Title       string     `json:"title,omitempty" db:"title" validate:"required"`
	Description string     `json:"description,omitempty" db:"description" validate:"required"`
	Released    CustomDate `json:"released" db:"released" validate:"required"`
	Rating      int8       `json:"rating,omitempty" db:"rating" validate:"required"`
}
