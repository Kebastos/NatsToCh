package models

import "time"

type DefaultTable struct {
	Subject        string
	CreateDateTime time.Time
	Content        string
}
