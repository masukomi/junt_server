package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type StatusChangesController struct {
	Db *gorm.DB
}

func (cc *StatusChangesController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	statusChange := models.StatusChange{}
	if err := r.DecodeJsonPayload(&statusChange); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&statusChange).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": statusChange.Id})

}

func (cc *StatusChangesController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	statusChange := models.StatusChange{}
	if cc.Db.First(&statusChange, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&statusChange)
}

func (cc *StatusChangesController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	statusChange := models.StatusChange{}
	if cc.Db.First(&statusChange, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := cc.Db.Delete(&statusChange).Error; err != nil {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
		return
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *StatusChangesController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	statusChanges := []models.StatusChange{}
	cc.Db.Find(&statusChanges)
	w.WriteJson(&statusChanges)
}
