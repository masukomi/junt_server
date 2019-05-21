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
	Db *gorm.DB
}

func (cc *JobsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	job := models.Job{}
	if err := r.DecodeJsonPayload(&job); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&job).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": job.Id})

}

func (cc *JobsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	job := models.Job{}
	if cc.Db.First(&job, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&job)
}

func (cc *JobsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	jobs := []models.Job{}
	cc.Db.Find(&jobs)
	w.WriteJson(&jobs)
}
