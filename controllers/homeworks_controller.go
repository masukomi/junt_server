package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type HomeworksController struct {
	CrudControllerImpl
	Db *gorm.DB
}

func (cc *HomeworksController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	homework := models.Homework{}
	if err := r.DecodeJsonPayload(&homework); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&homework).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": homework.Id})

}

func (cc *HomeworksController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	homework := models.Homework{}
	if cc.Db.First(&homework, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&homework)
}

func (cc *HomeworksController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	homework := models.Homework{}
	if cc.Db.First(&homework, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := homework.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *HomeworksController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	homeworks := []models.Homework{}
	cc.Db.Find(&homeworks)
	w.WriteJson(&homeworks)
}
