package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

var dateFormat = "2006-01-02"

type CustomDate struct {
	Date time.Time
}

func (d *CustomDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(fmt.Sprintf(`"%s"`, dateFormat), string(b))
	if err != nil {
		return err
	}
	d.Date = date
	return
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Date.Format(dateFormat))
}
