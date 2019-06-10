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
func GetPersonEvents(db *gorm.DB, personIds ...int64) ([]IEvent, error) {
	if len(personIds) > 1 {
		return []IEvent{}, errors.New("maximum of one person per request")
	}

	// TODO: figure out some way to
	// make this less... manual
	followups := []Followup{}
	interviews := []Interview{}
	thanksEmails := []ThanksEmail{}

	whereClause := GenerateIEventWhereClause(personIds...)
	// db.Where("job_id = ?", job_id).Find(&homeworks)
	db.Preload("People").Where(whereClause).Find(&followups)
	db.Preload("People").Where(whereClause).Find(&interviews)
	db.Preload("People").Where(whereClause).Find(&thanksEmails)

	size := len(followups) +
		len(interviews) +
		len(thanksEmails)
	iEvents := GroupRandomIEvents(size,
		followups,
		interviews,
		thanksEmails,
	)
	// sort them by CreatedAt
	sort.Sort(ByCreationDate(iEvents))
	return iEvents, nil

}
func GenerateIEventPersonWhereClauseause(personIds ...int64) string {
	if len(jobIds) > 0 {
		// person_id should only exist in the foo_people
		// and people_foo tables.
		return fmt.Sprintf("where  person_id = %v", personIds[0])
		// don't need to worry about SQL Injection because
		// it can't possibly be anything other than an int.
	} else {
		return "where true"
	}
}
