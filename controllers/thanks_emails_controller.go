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

	thanksEmail := models.ThanksEmail{}
	if err := r.DecodeJsonPayload(&thanksEmail); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := thanksEmail.ConvertIdsToPeople(cc.Db); err != nil {
		// TODO JSON ERROR STATUS
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cc.Db.Save(&thanksEmail).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": thanksEmail.Id})

}

func (cc *ThanksEmailsController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	thanksEmail := models.ThanksEmail{}
	if cc.Db.First(&thanksEmail, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&thanksEmail)
}

func (cc *ThanksEmailsController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	thanksEmail := models.ThanksEmail{}
	if cc.Db.First(&thanksEmail, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := thanksEmail.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *ThanksEmailsController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	thanksEmails := []models.ThanksEmail{}
	cc.Db.Find(&thanksEmails)
	w.WriteJson(&thanksEmails)
}

func (tec *ThanksEmailsController) Update(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	thanksEmail := models.ThanksEmail{}
	if tec.Db.First(&thanksEmail, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	// TODO decode into
	var data map[string]interface{}

	if err := r.DecodeJsonPayload(&data); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := thanksEmail.UpdateFromJson(data, tec.Db); err != nil {
		w.WriteJson(map[string]string{"status": "ERROR", "description": err.Error()})
		rest.NotFound(w, r)
		return

	}

	if err := tec.Db.Save(&thanksEmail).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}
