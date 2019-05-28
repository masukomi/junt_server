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
	Db *gorm.DB
}

func (cc *FollowupsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	followup := models.Followup{}
	if err := r.DecodeJsonPayload(&followup); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := followup.ConvertIdsToPeople(cc.Db); err != nil {
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
	if cc.Db.First(&followup, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&followup)
}

func (cc *FollowupsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	followups := []models.Followup{}
	cc.Db.Find(&followups)
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
