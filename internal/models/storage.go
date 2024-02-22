package models

import (
	"github.com/google/uuid"
	"time"
)

type DefaultEntity struct {
	Id             string
	ClientId       string
	Subject        string
	CreateDateTime time.Time
	Content        string
}

func NewDefaultEntity(clientName string, subject string, msg string) *DefaultEntity {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil
	}

	return &DefaultEntity{
		Id:             u.String(),
		ClientId:       clientName,
		Subject:        subject,
		CreateDateTime: time.Now(),
		Content:        msg,
	}
}
