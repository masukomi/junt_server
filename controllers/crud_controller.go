package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type CrudController interface {
	Create(w rest.ResponseWriter, r *rest.Request)
	// READ...
	ListAll(w rest.ResponseWriter, r *rest.Request)
	FindById(w rest.ResponseWriter, r *rest.Request)
	// TODO implement update
	// Update(w rest.ResponseWriter, r *rest.Request)
	Delete(w rest.ResponseWriter, r *rest.Request)
	UpdateModel(model models.IJsonUpdateable,
		db *gorm.DB,
		w rest.ResponseWriter,
		r *rest.Request)
}

type CrudControllerImpl struct {
}

func (cci *CrudControllerImpl) UpdateModel(model models.IJsonUpdateable,
	db *gorm.DB,
	w rest.ResponseWriter,
	r *rest.Request) {

	// TODO decode into
	var data map[string]interface{}

	if err := r.DecodeJsonPayload(&data); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := model.UpdateFromJson(data, db); err != nil {
		w.WriteJson(map[string]string{"status": "ERROR", "description": err.Error()})
		rest.Error(w, "JSON didn't meet API expectations", http.StatusUnprocessableEntity)
		return

	}

	if err := db.Save(&model).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}
