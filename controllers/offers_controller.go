package controllers

import (
	// "log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	"net/http"
)

type OffersController struct {
	CrudControllerImpl
	Db *gorm.DB
}

func (cc *OffersController) Create(w rest.ResponseWriter,
	r *rest.Request) {

	offer := models.Offer{}
	if err := r.DecodeJsonPayload(&offer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cc.Db.Save(&offer).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(
		map[string]interface{}{"status": "SUCCESS", "id": offer.Id})

}

func (cc *OffersController) FindById(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	offer := models.Offer{}
	if cc.Db.First(&offer, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&offer)
}
func (oc *OffersController) Update(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	offer := models.Offer{}
	if oc.Db.First(&offer, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	oc.UpdateModel(&offer, oc.Db, w, r)
	// see comment in UpdateModel for why this isn't there
	if err := oc.Db.Save(&offer).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "SUCCESS"})
}

func (cc *OffersController) Delete(w rest.ResponseWriter,
	r *rest.Request) {

	id := r.PathParam("id")
	offer := models.Offer{}
	if cc.Db.First(&offer, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	success, err := offer.HolisticDeletion(cc.Db)
	if success {
		w.WriteJson(map[string]string{"status": "SUCCESS"})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *OffersController) ListAll(w rest.ResponseWriter,
	r *rest.Request) {
	offers := []models.Offer{}
	cc.Db.Find(&offers)
	w.WriteJson(&offers)
}
