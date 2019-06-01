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
	CrudControllerImpl
	Db *gorm.DB
}

func (cc *InterviewsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	interview := models.Interview{}
	if err := r.DecodeJsonPayload(&interview); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := interview.ConvertIdsToPeople(cc.Db); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cc.Db.Save(&interview).Error; err != nil {
		// TODO JSON ERROR STATUS
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

func (cc *InterviewsController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	interview := models.Interview{}
	if cc.Db.First(&interview, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := interview.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (ic *InterviewsController) Update(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	interview := models.Interview{}
	if ic.Db.First(&interview, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	// TODO decode into
	var data map[string]interface{}

	if err := r.DecodeJsonPayload(&data); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := interview.UpdateFromJson(data, ic.Db); err != nil {
		w.WriteJson(map[string]string{"status": "ERROR", "description": err.Error()})
		rest.Error(w, "JSON didn't meet API expectations", http.StatusUnprocessableEntity)
		return

	}

	if err := ic.Db.Save(&interview).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}
