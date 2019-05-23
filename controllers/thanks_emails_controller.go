package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type ThanksEmailsController struct {
	Db *gorm.DB
}

func (cc *ThanksEmailsController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	thanks_email := models.ThanksEmail{}
	if err := r.DecodeJsonPayload(&thanks_email); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&thanks_email).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": thanks_email.Id})

}

func (cc *ThanksEmailsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	thanks_email := models.ThanksEmail{}
	if cc.Db.First(&thanks_email, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&thanks_email)
}

func (cc *ThanksEmailsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	thanks_emails := []models.ThanksEmail{}
	cc.Db.Find(&thanks_emails)
	w.WriteJson(&thanks_emails)
}
