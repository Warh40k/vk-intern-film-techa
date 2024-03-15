package domain

import (
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

type CustomDate time.Time

func (d *CustomDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(fmt.Sprintf(`"%s"`, dateFormat), string(b))
	if err != nil {
		return err
	}
	*d = CustomDate(date)
	return
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (ct CustomDate) String() string {
	t := time.Time(ct)
	return fmt.Sprintf("%q", t.Format(dateFormat))
}
