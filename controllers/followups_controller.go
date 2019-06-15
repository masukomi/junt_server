package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type FollowupsController struct {
	CrudControllerImpl
	Db *gorm.DB
}

func (cc *FollowupsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	followup := models.Followup{}
	if err := r.DecodeJsonPayload(&followup); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followup.ConvertIdToJob(cc.Db)
	// all events have a Job but not all events have People
	if err := models.ConvertIdsToPeople(cc.Db, &followup); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cc.Db.Save(&followup).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": followup.Id})

}

func (cc *FollowupsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	followup := models.Followup{}
	if cc.Db.Preload("People").First(&followup, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&followup)
}
func (cc *FollowupsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	followups := []models.Followup{}
	cc.Db.Preload("People").Find(&followups)
	w.WriteJson(&followups)
}

func (cc *FollowupsController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	followup := models.Followup{}
	if cc.Db.First(&followup, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := followup.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (fc *FollowupsController) Update(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	followup := models.Followup{}
	if fc.Db.First(&followup, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	// // TODO decode into
	// var data map[string]interface{}
	//
	// if err := r.DecodeJsonPayload(&data); err != nil {
	// 	rest.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if err := followup.UpdateFromJson(data, fc.Db); err != nil {
	// 	w.WriteJson(map[string]string{"status": "ERROR", "description": err.Error()})
	// 	rest.Error(w, "JSON didn't meet API expectations", http.StatusUnprocessableEntity)
	// 	return
	//
	// }

	fc.UpdateModel(&followup, fc.Db, w, r)
	// see comment in UpdateModel for why this isn't there
	if err := fc.Db.Save(&followup).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}
