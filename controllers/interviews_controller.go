package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type InterviewsController struct {
	Db *gorm.DB
}

func (cc *InterviewsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	interview := models.Interview{}
	if err := r.DecodeJsonPayload(&interview); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&interview).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": interview.Id})

}

func (cc *InterviewsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	interview := models.Interview{}
	if cc.Db.First(&interview, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&interview)
}

func (cc *InterviewsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	interviews := []models.Interview{}
	cc.Db.Find(&interviews)
	w.WriteJson(&interviews)
}
