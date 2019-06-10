package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Followup struct {
	Event
	People    []Person `json:"-" gorm:"many2many:followups_people;"`
	PersonIds []int64  `json:"person_ids" gorm:"-"`
}

func (f *Followup) GetPersonIds() []int64 {
	return f.PersonIds
}
func (f *Followup) GetPeople() []Person {
	return f.People
}
func (f *Followup) ExtractPeople(db *gorm.DB) []Person {
	ConvertIdsToPeople(db, f)
	return f.People
}
func (f *Followup) SetPeople(peeps []Person) {
	f.People = peeps
}
func (f *Followup) SetPersonIds(personIds []int64) {
	f.PersonIds = personIds
}

func (f *Followup) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return UpdatePeopleEventFromJson(data, db, f)
}
func (f *Followup) MarshalJSON() ([]byte, error) {
	f.PersonIds = ExtractPersonIds(f)

	type Alias Followup
	return json.Marshal(&struct {
		EventType string `json:"event_type" gorm:"-"`
		*Alias
	}{
		Alias:     (*Alias)(f),
		EventType: "followup",
	})
}
