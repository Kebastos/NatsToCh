package models

import "time"

type DefaultTable struct {
	Id             string
	Subject        string
	CreateDateTime time.Time
	Content        string
}
