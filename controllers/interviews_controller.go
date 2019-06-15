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

func (ic *InterviewsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	interview := models.Interview{}
	if err := r.DecodeJsonPayload(&interview); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := models.ConvertIdsToPeople(ic.Db, &interview); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ic.Db.Save(&interview).Error; err != nil {
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
	if cc.Db.Preload("People").First(&interview, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&interview)
}

func (cc *InterviewsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	interviews := []models.Interview{}
	cc.Db.Preload("People").Find(&interviews)
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

	ic.UpdateModel(&interview, ic.Db, w, r)
	// see comment in UpdateModel for why this isn't there
	if err := ic.Db.Save(&interview).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}
