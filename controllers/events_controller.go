package controllers

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"sort"
	"strconv"
)

type EventsController struct {
	Db *gorm.DB
}

func (ec *EventsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	events, _ := ec.getIEvents() // could pass in a job id
	// err is only if you have > 1 jobId which we don't ;)
	w.WriteJson(&events)
}

func (ec *EventsController) ListAllForJob(w rest.ResponseWriter,
	r *rest.Request) {
	strId := r.PathParam("id")
	id, err := strconv.Atoi(strId)
	if err == nil {
		events, _ := ec.getIEvents(id)
		// err is only if you have > 1 jobId which we don't ;)
		w.WriteJson(&events)
	} else {
		w.WriteJson(map[string]string{"status": "ERROR",
			"description": "no job id found: " + err.Error()})
	}
}

func (ec *EventsController) getIEvents(jobIds ...int) ([]models.IEvent, error) {
	if len(jobIds) > 1 {
		return []models.IEvent{}, errors.New("maximum of one job per request")
	}

	// TODO: figure out some way to
	// make this less... manual
	followups := []models.Followup{}
	homeworks := []models.Homework{}
	interviews := []models.Interview{}
	offers := []models.Offer{}
	statusChanges := []models.StatusChange{}
	thanksEmails := []models.ThanksEmail{}

	whereClause := ec.generateWhereClause(jobIds...)
	// ec.Db.Where("job_id = ?", job_id).Find(&homeworks)
	ec.Db.Where(whereClause).Find(&followups)
	ec.Db.Where(whereClause).Find(&homeworks)
	ec.Db.Where(whereClause).Find(&interviews)
	ec.Db.Where(whereClause).Find(&offers)
	ec.Db.Where(whereClause).Find(&statusChanges)
	ec.Db.Where(whereClause).Find(&thanksEmails)
	ec.Db.Where(whereClause).Find(&thanksEmails)

	iEvents := ec.groupRandomIEvents(followups,
		homeworks,
		interviews,
		offers,
		statusChanges,
		thanksEmails,
	)
	// sort them by CreatedAt
	sort.Sort(models.ByCreationDate(iEvents))
	return iEvents, nil
}
func (ec *EventsController) generateWhereClause(jobIds ...int) string {
	if len(jobIds) > 0 {
		return fmt.Sprintf("where job_id = %v", jobIds[0])
		// don't need to worry about SQL Injection because
		// it can't possibly be anything other than an int.
	} else {
		return "where true"
	}
}

func (ec *EventsController) groupRandomIEvents(iEvents ...interface{}) []models.IEvent {
	size := 0
	for _, x := range iEvents {
		size += len(x.([]models.IEvent))
	}
	response := make([]models.IEvent, size)
	idx := 0
	for _, currentSlice := range iEvents {
		// each element of iEvents is itself a slice
		for _, event := range currentSlice.([]models.IEvent) {
			response[idx] = event //.(models.IEvent)
			idx++
		}
	}

	return response
}
