package domain

type Film struct {
	Id          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title" validate:"required,gt=0,lte=150"`
	Description string     `json:"description" db:"description" validate:"required,lte=1000"`
	Released    CustomDate `json:"released" db:"released" validate:"required"`
	Rating      int8       `json:"rating" db:"rating" validate:"gte=0,lte=10"`
}

type NullableFilm struct {
	Id          int         `json:"-"`
	Title       *string     `json:"title" db:"title" validate:"omitempty,gt=0,lte=150"`
	Description *string     `json:"description" db:"description" validate:"omitempty,lte=1000"`
	Released    *CustomDate `json:"released" db:"released" validate:"omitempty"`
	Rating      *int8       `json:"rating" db:"rating" validate:"omitempty,gte=0,lte=10"`
	ActorIds    []int       `json:"actorIds" db:"-"`
}
