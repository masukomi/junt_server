package models

import (
	"time"
)

type Interview struct {
	Event
	ScheduledAt time.Time `json:"scheduled_at"`
	// generated if not supplied
	Length int64 `json:"length"`
	// in minutes
	Rating string `sql:"type:text;" json:"rating"`
	// emoji
	Type string `sql:"type:text;" json:"type"`
	// string enum (user created)
	People []Person `gorm:"many2many:interviews_people;"`
	// has and belongs to many jobs
}
