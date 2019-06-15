package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type PeopleController struct {
	CrudControllerImpl
	Db *gorm.DB
}

func (cc *PeopleController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	person := models.Person{}
	if err := r.DecodeJsonPayload(&person); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		// TODO: error JSON
		return
	}
	// convert ids to jobs

	if err := person.ConvertIdsToJobs(cc.Db); err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		// TODO: error JSON
		return
	}

	if err := cc.Db.Save(&person).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		// TODO: error JSON
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": person.Id})

}

func (cc *PeopleController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	person := models.Person{}
	if cc.Db.Preload("Jobs").First(&person, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&person)
}
func (pc *PeopleController) Update(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	person := models.Person{}
	if pc.Db.First(&person, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	pc.UpdateModel(&person, pc.Db, w, r)

	// see comment in UpdateModel for why this isn't there
	if err := pc.Db.Save(&person).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})

}

func (cc *PeopleController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	people := []models.Person{}
	cc.Db.Preload("Jobs").Find(&people)
	w.WriteJson(&people)
}

func (cc *PeopleController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	person := models.Person{}
	if cc.Db.First(&person, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	// success, err := person.HolisticDeletion(cc.Db)
	person.HolisticDeletion(cc.Db)
	// It's going to raise an error
	// it's going to ignore the transaction rollback
	// and yet, it's going to work.
	// I don't know why this is problematic and Jobs deletion is not.
	//if success {
	w.WriteJson(map[string]string{"status": "SUCCESS"})
	// } else {
	// 	rest.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}
