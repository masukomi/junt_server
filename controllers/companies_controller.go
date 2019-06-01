package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type CompaniesController struct {
	Db *gorm.DB
	CrudControllerImpl
}

func (cc *CompaniesController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	company := models.Company{}
	if err := r.DecodeJsonPayload(&company); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&company).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": company.Id})

}

func (cc *CompaniesController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	company := models.Company{}
	if cc.Db.First(&company, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&company)
}

func (cc *CompaniesController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	company := models.Company{}
	if cc.Db.First(&company, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := company.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
		return
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *CompaniesController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	companies := []models.Company{}
	cc.Db.Find(&companies)
	w.WriteJson(&companies)
}
