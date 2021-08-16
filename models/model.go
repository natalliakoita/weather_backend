package models

import "time"

type Weather struct {
	ID          int
	City        string
	TimeStamp   time.Time
	Temperature float32
}
