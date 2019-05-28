package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type CrudController interface {
	Create(w rest.ResponseWriter, r *rest.Request)
	// READ...
	ListAll(w rest.ResponseWriter, r *rest.Request)
	FindById(w rest.ResponseWriter, r *rest.Request)
	Update(w rest.ResponseWriter, r *rest.Request)
	Delete(w rest.ResponseWriter, r *rest.Request)
}
