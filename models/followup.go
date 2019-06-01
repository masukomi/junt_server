package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Followup struct {
	PeopleEvent
	People []Person `gorm:"many2many:followups_people;"`
}

func (f *Followup) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return f.UpdatePeopleEventFromJson(data, db)
}
