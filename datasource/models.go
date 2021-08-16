package datasource

import "time"

type Weather struct {
	ID          int       `db:"id"`
	City        string    `db:"city"`
	TimeStamp   time.Time `db:"dt"`
	Temperature float32   `db:"temperature"`
}
