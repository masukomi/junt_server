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

	// id -- ignore
	// created_at time.Time
	// updated_at time.Time
	// note string (markdown)
	// job_id -- never expect this to change but will allow
	// person_ids -- will be converted to Person objects
	//
	for key, value := range data {
		switch key {
		case "created_at":
			newTime, err := time.Parse(time.RFC3339, value.(string))
			if err != nil {
				te.CreatedAt = newTime
			} else {
				return errors.New("invalid created_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		case "updated_at":
			newTime, err := time.Parse(time.RFC3339, value.(string))
			if err != nil {
				te.UpdatedAt = newTime
			} else {
				return errors.New("invalid created_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		case "note":
			te.Note = value.(string)
		case "job_id":
			te.JobId = int64(value.(float64))
			job := Job{}
			if db.First(&job, value).Error != nil {
				return errors.New("invalid associated job_id: " + strconv.FormatInt(te.JobId, 10))
				// why can I convert from an int with Itoa but not an int64?
			}
			te.Job = job
		case "person_ids":
			person_ids := []int64{}
			for _, num := range value.([]interface{}) { // []interface{}
				person_ids = append(person_ids, int64(num.(float64)))
			}
			te.PersonIds = person_ids
			if err := te.ConvertIdsToPeople(db); err != nil {
				return errors.New("invalid associated person_ids")
			}
		}
	}
	return nil
}
