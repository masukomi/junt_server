package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Offer struct {
	Event
	Status string `sql:"type:text" json:"status"`
}

func (o *Offer) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return o.UpdateEventFromJson(data, db)
}
