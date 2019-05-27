package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"sort"
	"time"
)

type IEvent interface {
	CreationDate() time.Time
	HolisticDeletion(db *gorm.DB) (bool, error)
}

// the folowing BS is because Go hates you
// and doesn't want you to have anything nice
// like Ruby's Comparable functionality.
// I'm sorry. Really.
type ByCreationDate []IEvent

func (e ByCreationDate) Len() int {
	return len(e)
}

func (e ByCreationDate) Less(i, j int) bool {
	return e[i].CreationDate().Unix() < e[j].CreationDate().Unix()
}

func (e ByCreationDate) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

/// end madness... for now

func GetIEvents(db *gorm.DB, jobIds ...int) ([]IEvent, error) {
	if len(jobIds) > 1 {
		return []IEvent{}, errors.New("maximum of one job per request")
	}

	// TODO: figure out some way to
	// make this less... manual
	followups := []Followup{}
	homeworks := []Homework{}
	interviews := []Interview{}
	offers := []Offer{}
	statusChanges := []StatusChange{}
	thanksEmails := []ThanksEmail{}

	whereClause := GenerateIEventWhereClause(jobIds...)
	// db.Where("job_id = ?", job_id).Find(&homeworks)
	db.Where(whereClause).Find(&followups)
	db.Where(whereClause).Find(&homeworks)
	db.Where(whereClause).Find(&interviews)
	db.Where(whereClause).Find(&offers)
	db.Where(whereClause).Find(&statusChanges)
	db.Where(whereClause).Find(&thanksEmails)
	db.Where(whereClause).Find(&thanksEmails)

	iEvents := GroupRandomIEvents(followups,
		homeworks,
		interviews,
		offers,
		statusChanges,
		thanksEmails,
	)
	// sort them by CreatedAt
	sort.Sort(ByCreationDate(iEvents))
	return iEvents, nil
}
func GenerateIEventWhereClause(jobIds ...int) string {
	if len(jobIds) > 0 {
		return fmt.Sprintf("where job_id = %v", jobIds[0])
		// don't need to worry about SQL Injection because
		// it can't possibly be anything other than an int.
	} else {
		return "where true"
	}
}

func GroupRandomIEvents(iEvents ...interface{}) []IEvent {
	size := 0
	for _, x := range iEvents {
		size += len(x.([]IEvent))
	}
	response := make([]IEvent, size)
	idx := 0
	for _, currentSlice := range iEvents {
		// each element of iEvents is itself a slice
		for _, event := range currentSlice.([]IEvent) {
			response[idx] = event //.(IEvent)
			idx++
		}
	}

	return response
}
