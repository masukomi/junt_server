package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PeopleEvent struct {
	Event
	People    []Person
	PersonIds []int64 `gorm:"-"`
}

func (pe *PeopleEvent) GetPersonIds() []int64 {

	person_ids := make([]int64, len(pe.People))
	for i, p := range pe.People {
		person_ids[i] = p.Id
	}
	return person_ids
}

/// JSON RELATED STUFF
// NOTE: unmarshall via normal means
// THEN call ConvertIdsToPeople

func (pe *PeopleEvent) MarshalJSON() ([]byte, error) {
	person_ids := pe.GetPersonIds()
	type Alias PeopleEvent
	return json.Marshal(&struct {
		PersonIds []int64 `json:"person_ids"`
		*Alias
	}{
		Alias:     (*Alias)(pe),
		PersonIds: person_ids,
	})
}
func (pe *PeopleEvent) SetPeople(peeps []Person) {
	pe.People = peeps
}
func (pe *PeopleEvent) ConvertIdsToPeople(db *gorm.DB) error {
	peeps := make([]Person, len(pe.GetPersonIds()))
	for i, pid := range pe.GetPersonIds() {
		person := Person{}
		if err := db.First(&person, pid).Error; err != nil {
			return err
		}
		peeps[i] = person
	}
	pe.SetPeople(peeps)
	return nil
}
