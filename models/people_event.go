package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sort"
)

type IPersonEvent interface {
	GetPeople() []Person                // gets People from People
	ExtractPeople(db *gorm.DB) []Person // gets people via PersonIds
	GetPersonIds() []int64              // gets the ids from PersonIds
	// ExtractPersonIds() []int64 // gets the ids via People
	SetPeople(peeps []Person)
	SetPersonIds(personIds []int64)
	UpdateEventFromJson(data map[string]interface{}, db *gorm.DB) error
}

// type PeopleEvent struct {
// 	Event
// 	People    []Person `json:"-"`
// 	PersonIds []int64  `json:"person_ids" gorm:"-"`
// }

func ExtractPersonIds(ipe IPersonEvent) []int64 {
	peeps := ipe.GetPeople()
	person_ids := make([]int64, len(peeps))
	for i, p := range peeps {
		person_ids[i] = p.Id
	}
	return person_ids
}

func ConvertIdsToPeople(db *gorm.DB, ipe IPersonEvent) error {
	peep_ids := ipe.GetPersonIds()
	peeps := make([]Person, len(peep_ids))
	for i, pid := range peep_ids {
		person := Person{}
		if err := db.First(&person, pid).Error; err != nil {
			return err
		}
		peeps[i] = person
	}
	ipe.SetPeople(peeps)
	return nil
}
func UpdatePeopleEventFromJson(data map[string]interface{}, db *gorm.DB, ipe IPersonEvent) error {
	if err := ipe.UpdateEventFromJson(data, db); err != nil {
		return err
	}
	value, ok := data["personIds"]
	if ok {
		personIds := ExtractIdsFromJsonArray(value.([]interface{}))
		ipe.SetPersonIds(personIds)
		if err := ipe.ExtractPeople(db); err != nil {
			return errors.New("invalid associated personIds")
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
func GenerateIEventPersonWhereClause(personIds ...int64) string {
	if len(personIds) > 0 {
		// person_id should only exist in the foo_people
		// and people_foo tables.
		return fmt.Sprintf("where  person_id = %v", personIds[0])
		// don't need to worry about SQL Injection because
		// it can't possibly be anything other than an int.
	} else {
		return "where true"
	}
}
