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
	Event
	ScheduledAt time.Time `json:"scheduled_at"`
	// generated if not supplied
	Length int64 `json:"length"`
	// in minutes
	Rating string `sql:"type:text;" json:"rating"`
	// emoji
	Type string `sql:"type:text;" json:"type"`
	// string enum (user created)
	People []Person `json:"-" gorm:"many2many:interviews_people;"`
	// has and belongs to many people
	PersonIds []int64 `json:"person_ids" gorm:"-"`
}

func (i *Interview) GetPersonIds() []int64 {
	return i.PersonIds
}
func (i *Interview) GetPeople() []Person {
	return i.People
}
func (i *Interview) SetPeople(peeps []Person) {
	i.People = peeps
}
func (i *Interview) SetPersonIds(personIds []int64) {
	i.PersonIds = personIds
}
func (i *Interview) ExtractPeople(db *gorm.DB) []Person {
	ConvertIdsToPeople(db, i)
	return i.People
}

func (i *Interview) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	if err := UpdatePeopleEventFromJson(data, db, i); err != nil {
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
	i.PersonIds = ExtractPersonIds(i)
	fmt.Println("interview: ", i.Id, " person_ids: ", i.PersonIds)
	type Alias Interview
	return json.Marshal(&struct {
		EventType string `json:"event_type" gorm:"-"`
		*Alias
	}{
		Alias:     (*Alias)(i),
		EventType: "interview",
	})
}
