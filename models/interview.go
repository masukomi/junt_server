package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Interview struct {
	PeopleEvent
	ScheduledAt time.Time `json:"scheduled_at"`
	// generated if not supplied
	Length int64 `json:"length"`
	// in minutes
	Rating string `sql:"type:text;" json:"rating"`
	// emoji
	Type string `sql:"type:text;" json:"type"`
	// string enum (user created)
	People []Person `json:"-" gorm:"many2many:interviews_people;"`
	// has and belongs to many jobs
}

func (i *Interview) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	if err := i.UpdatePeopleEventFromJson(data, db); err != nil {
		return err
	}
	// also has...
	// scheduled_at: date_time
	// length:       int
	// 			  (minutes)
	// rating:       string enum
	// 			  (emoji suggested)
	// type:         string enum
	for key, value := range data {
		switch key {
		case "scheduled_at":
			newTime, err := time.Parse(time.RFC3339, value.(string))
			if err != nil {
				i.ScheduledAt = newTime
			} else {
				return errors.New("invalid scheduled_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		case "length":
			i.Length = int64(value.(float64))
		case "rating":
			i.Rating = value.(string)
		case "type":
			i.Type = value.(string)
		}
	}

	return nil
}
func (i *Interview) MarshalJSON() ([]byte, error) {
	personIds := i.GetPersonIds()
	type Alias Interview
	return json.Marshal(&struct {
		PersonIds []int64 `json:"person_ids"`
		*Alias
	}{
		Alias:     (*Alias)(i),
		PersonIds: personIds,
	})
}
