package models

type Offer struct {
	Event
	Status string `sql:"type:text" json:"status"`
}
