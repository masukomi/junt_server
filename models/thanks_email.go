package models

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ThanksEmail struct {
	Event
	People    []Person `json:"-" gorm:"many2many:people_thanks_emails;"` // has and belongs to many people
	PersonIds []int64  `json:"person_ids" gorm:"-"`
}

func (te *ThanksEmail) GetPersonIds() []int64 {
	return te.PersonIds
}
func (te *ThanksEmail) GetPeople() []Person {
	return te.People
}
func (te *ThanksEmail) SetPeople(peeps []Person) {
	te.People = peeps
}
func (te *ThanksEmail) SetPersonIds(personIds []int64) {
	te.PersonIds = personIds
}
func (te *ThanksEmail) ExtractPeople(db *gorm.DB) []Person {
	ConvertIdsToPeople(db, te)
	return te.People
}

func (te *ThanksEmail) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return UpdatePeopleEventFromJson(data, db, te)
}
func (te *ThanksEmail) MarshalJSON() ([]byte, error) {
	fmt.Println("XXX te id: ", te.Id, " people: ", te.People, "person_ids: ", te.PersonIds)
	te.PersonIds = ExtractPersonIds(te)
	type Alias ThanksEmail
	return json.Marshal(&struct {
		EventType string `json:"event_type" gorm:"-"`
		*Alias
	}{
		Alias:     (*Alias)(te),
		EventType: "thanks_email",
	})
}
