package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PeopleEvent struct {
	Event
	People    []Person `json:"-"`
	PersonIds []int64  `json:"person_ids" gorm:"-"`
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

// func (pe *PeopleEvent) MarshalJSON() ([]byte, error) {
// 	fmt.Println("XXXXX in PeopleEvent#MarshalJSON")
// 	personIds := pe.GetPersonIds()
// 	type Alias PeopleEvent
// 	return json.Marshal(&struct {
// 		PersonIds []int64 `json:"person_ids"`
// 		*Alias
// 	}{
// 		Alias:     (*Alias)(pe),
// 		PersonIds: personIds,
// 	})
// }
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
func (pe *PeopleEvent) UpdatePeopleEventFromJson(data map[string]interface{}, db *gorm.DB) error {
	if err := pe.UpdateEventFromJson(data, db); err != nil {
		return err
	}
	value, ok := data["person_ids"]
	if ok {
		person_ids := []int64{}
		for _, num := range value.([]interface{}) { // []interface{}
			person_ids = append(person_ids, int64(num.(float64)))
		}
		pe.PersonIds = person_ids
		if err := pe.ConvertIdsToPeople(db); err != nil {
			return errors.New("invalid associated person_ids")
		}
	}
	// if not ok, no worries. they weren't updating that association
	return nil
}
