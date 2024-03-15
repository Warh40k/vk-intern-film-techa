package domain

type Actor struct {
	Id       int        `json:"id,omitempty" db:"id"`
	Name     string     `json:"name,omitempty" db:"name"`
	Gender   int        `json:"gender,omitempty" db:"gender"`
	Birthday CustomDate `json:"birthday" db:"birthday"`
}
