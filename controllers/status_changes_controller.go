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

	status_change := models.StatusChange{}
	if err := r.DecodeJsonPayload(&status_change); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&status_change).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": status_change.Id})

}

func (cc *StatusChangesController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	status_change := models.StatusChange{}
	if cc.Db.First(&status_change, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&status_change)
}

func (cc *StatusChangesController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	status_changes := []models.StatusChange{}
	cc.Db.Find(&status_changes)
	w.WriteJson(&status_changes)
}
