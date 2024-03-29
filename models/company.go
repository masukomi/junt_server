package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

/*
Company
----------------
id:         int
name:       string
location:   string (optional)
url:        string (url optional)
note:       string (optional - markdown)
created_at: date_time (optional)
updated_at: date_time (optional)

has_many:   jobs (maybe)
has_many:   people (through applications)
*/

type Company struct {
	// gorm.Model
	// Id, CreatedAt, UpdatedAt would normally come from gorm.Model but we need to specify
	// the json keys for them so...
	Id        int64     `json:"id" gorm:"PRIMARY_KEY";"AUTO_INCREMENT" `         // generated by DB
	Name      string    `gorm:"index:co_names_idx" sql:"type:text;" json:"name"` // only required field
	Location  string    `sql:"type:text;" json:"location"`
	Url       string    `sql:"type:text;" json:"url"`
	Note      string    `sql:"type:text;" json:"note"` // markdown
	CreatedAt time.Time `json:"created_at"`            // generated if not supplied
	UpdatedAt time.Time `json:"updated_at"`            // generated if not supplied
	Jobs      []Job     `json:"-" gorm:"foreignkey:CompanyId"`
	People    []Person  `json:"-" gorm:"foreignkey:CompanyId"`
	Db        *gorm.DB  `json:"-" gorm:"-"`
}

func (c Company) HolisticDeletion(db *gorm.DB) (bool, error) {

	transaction := db.Begin()
	// find all the jobs
	for _, job := range c.Jobs {
		// returns slice of IEvent objects
		// tied to the job
		success, err := job.TransactionlessHolisticDeletion(db)
		if !success {
			return success, err
		}
	}
	for _, person := range c.People {
		success, err := DeleteOrRollback(transaction, db, person)
		if !success {
			return success, err
		}
	}
	// delete the company
	if err := db.Delete(&c).Error; err != nil {
		transaction.Commit()
		return true, nil
	} else {
		transaction.Rollback()
		return false, err
	}
}

func (c *Company) MarshalJSON() ([]byte, error) {
	jobs := c.Jobs
	jobIds := make([]int64, len(jobs))
	for idx, job := range jobs {
		jobIds[idx] = job.Id
	}

	people := c.People
	personIds := make([]int64, len(people))
	for idx, person := range people {
		personIds[idx] = person.Id
	}

	type Alias Company
	return json.Marshal(&struct {
		JobIds    []int64 `json:"job_ids"`
		PersonIds []int64 `json:"person_ids"`
		*Alias
	}{
		JobIds:    jobIds,
		PersonIds: personIds,
		Alias:     (*Alias)(c),
	})
}

func (c *Company) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	for key, value := range data {
		switch key {
		case "created_at":
			mTime := MaybeTimeFromValue(value.(string))
			if !mTime.IsError() {
				c.CreatedAt = mTime.Just
			} else {
				return errors.New("invalid created_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		case "updated_at":
			mTime := MaybeTimeFromValue(value.(string))
			if !mTime.IsError() {
				c.UpdatedAt = mTime.Just
			} else {
				return errors.New("invalid created_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		case "note":
			c.Note = value.(string)
		case "name":
			if value == nil {
				return errors.New("Companies must have a name")
			}
			c.Name = value.(string)
		case "url":
			c.Url = value.(string)
		case "location":
			c.Location = value.(string)
		}
	}
	return nil
}
