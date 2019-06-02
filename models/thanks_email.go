package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ThanksEmail struct {
	PeopleEvent
	People []Person `gorm:"many2many:people_thanks_emails;"` // has and belongs to many jobs
}

func (te *ThanksEmail) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	return te.UpdatePeopleEventFromJson(data, db)
}
func (te *ThanksEmail) MarshalJSON() ([]byte, error) {
	personIds := te.GetPersonIds()
	type Alias ThanksEmail
	return json.Marshal(&struct {
		PersonIds []int64 `json:"person_ids"`
		*Alias
	}{
		Alias:     (*Alias)(te),
		PersonIds: personIds,
	})
}
