package controllers

import (
	// "errors"
	// "fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	// "sort"
	"strconv"
)

type EventsController struct {
	Db *gorm.DB
}

func (ec *EventsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	events, _ := models.GetIEvents(ec.Db) // could pass in a job id
	// err is only if you have > 1 jobId which we don't ;)
	w.WriteJson(&events)
}

func (ec *EventsController) ListAllForJob(w rest.ResponseWriter,
	r *rest.Request) {
	strId := r.PathParam("id")
	id, err := strconv.Atoi(strId)
	if err == nil {
		events, _ := models.GetIEvents(ec.Db, id)
		// err is only if you have > 1 jobId which we don't ;)
		w.WriteJson(&events)
	} else {
		w.WriteJson(map[string]string{"status": "ERROR",
			"description": "no job id found: " + err.Error()})
	}
}
