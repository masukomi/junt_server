package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type StatusChange struct {
	Event
	From string `sql:"type:text" json:"from"`
	To   string `sql:"type:text" json:"to"`
}

func (sc *StatusChange) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	err := sc.UpdateEventFromJson(data, db)
	if err != nil {
		if value, ok := data["from"]; ok {
			val := value.(string)
			if val != "" {
				sc.From = val
			} else {
				return errors.New("StatusEvent's \"from\" field can't be empty")
			}
		} else {
			return errors.New("StatusEvents must have a valid \"from\" field")
		}

	}
	if err != nil {
		if value, ok := data["to"]; ok {
			val := value.(string)
			if val != "" {
				sc.To = val
			} else {
				return errors.New("StatusEvent's \"to\" field can't be empty")
			}
		} else {
			return errors.New("StatusEvents must have a valid \"to\" field")
		}

	}
	return nil
}
