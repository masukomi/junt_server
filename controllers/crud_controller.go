package controllers

type CrudController interface {
	Create()
	// READ...
	listAll()
	FindById()
	Update()
	Delete()
}
