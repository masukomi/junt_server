package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type PeopleController struct {
	Db *gorm.DB
}

func (cc *PeopleController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	person := models.Person{}
	if err := r.DecodeJsonPayload(&person); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&person).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": person.Id})

}

func (cc *PeopleController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	person := models.Person{}
	if cc.Db.First(&person, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&person)
}

func (cc *PeopleController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	people := []models.Person{}
	cc.Db.Find(&people)
	w.WriteJson(&people)
}
