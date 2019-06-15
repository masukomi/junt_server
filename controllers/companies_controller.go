package controllers

import (
	// "log"
	"errors"
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

func (cc *CompaniesController) LoadCompanyFromRequest(r *rest.Request) models.MaybeCompany {
	id := r.PathParam("id")
	company := models.Company{}
	if cc.Db.Preload("Jobs").Preload("People").First(&company, id).Error != nil {
		return models.ErrorCompany(errors.New("can't find company with that id"))
	}
	return models.JustCompany(company)

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

	// id := r.PathParam("id")
	// company := models.Company{}
	// if cc.Db.Preload("Jobs").Preload("People").First(&company, id).Error != nil {
	// 	rest.NotFound(w, r)
	// 	return
	// }
	maybeCompany := cc.LoadCompanyFromRequest(r)

	if !maybeCompany.IsError() {
		w.WriteJson(&maybeCompany.Just)
	} else {
		id := r.PathParam("id")
		rest.NotFound(w, r)
		w.WriteJson(
			map[string]interface{}{"status": "ERROR",
				"description": ("Unable to find company with id" + id)})
	}
}

func (cc *CompaniesController) Update(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	company := models.Company{}
	if cc.Db.First(&company, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	cc.UpdateModel(&company, cc.Db, w, r)

	// see comment in UpdateModel for why this isn't there
	if err := cc.Db.Save(&company).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})

}

func (cc *CompaniesController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	company := models.Company{}
	if cc.Db.Preload("Jobs").Preload("People").First(&company, id).Error != nil {
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
	cc.Db.Preload("Jobs").Preload("People").Find(&companies)
	w.WriteJson(&companies)
}
