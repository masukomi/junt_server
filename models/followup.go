package models

import (
	"encoding/json"
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
func (f *Followup) MarshalJSON() ([]byte, error) {
	personIds := f.GetPersonIds()
	type Alias Followup
	return json.Marshal(&struct {
		PersonIds []int64 `json:"person_ids"`
		*Alias
	}{
		Alias:     (*Alias)(f),
		PersonIds: personIds,
	})
}
