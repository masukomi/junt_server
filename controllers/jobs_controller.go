package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type JobsController struct {
	CrudControllerImpl
	Db *gorm.DB
}

func (jc *JobsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	job := models.Job{}
	if err := r.DecodeJsonPayload(&job); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := jc.Db.Save(&job).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": job.Id})

}

func (jc *JobsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	job := models.Job{}
	if jc.Db.Preload("People").First(&job, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&job)
}

func (jc *JobsController) Edit(w rest.ResponseWriter,
	r *rest.Request) {
	id := r.PathParam("id")
	job := models.Job{}
	if jc.Db.First(&job, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := r.DecodeJsonPayload(&job); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := jc.Db.Save(&job).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": job.Id})
}
func (jc *JobsController) Update(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	job := models.Job{}
	if jc.Db.First(&job, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	jc.UpdateModel(&job, jc.Db, w, r)
	// see comment in UpdateModel for why this isn't there
	if err := jc.Db.Save(&job).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}

func (jc *JobsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	jobs := []models.Job{}
	jc.Db.Preload("People").Find(&jobs)
	w.WriteJson(&jobs)
}

func (jc *JobsController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	job := models.Job{}
	if jc.Db.First(&job, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := job.HolisticDeletion(jc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		w.WriteJson(map[string]string{"status": "ERROR",
			"description": err.Error()})
	}
}
