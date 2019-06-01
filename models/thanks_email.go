package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strconv"
	"time"
)

type ThanksEmail struct {
	PeopleEvent
	People []Person `gorm:"many2many:people_thanks_emails;"` // has and belongs to many jobs
}

func (te *ThanksEmail) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return te.UpdatePeopleEventFromJson(data, db)
}
