package database

import "time"

type Entry struct {
	Id          int64
	LogTime     time.Time
	Description string
}
