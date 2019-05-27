package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// ┌──────────────────────────────────────────────┐
// │Job                                           │
// │-----------------------------------------     │
// │id:                 int                       │
// │company_id:         int                       │
// │job_title:          string (optional)         │
// │posting_url:        string (url)              │
// │source:             string (optional)         │
// │referred_by:        string (optional)         │
// │salary_range:       string (optional)         │
// │application_method: string (optional)         │
// │note:               string (optional)         │
// │created_at:         date_time (optional)      │
// │updated_at:         date_time (optional)      │
// │start_date:         date_time (optional)      │
// │                                              │
// │                                              │
// │belongs_to:         company                   │
// │has_many:           events (maybe)            │
// │has_and_belongs_to_many:    people (hr contact(s))    │
// └──────────────────────────────────────────────┘
// dates are ISO 8601 RFC3339
// t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

type Job struct {
	// gorm.Model
	// Id, CreatedAt, UpdatedAt would normally come from gorm.Model but we need to specify
	// the json keys for them so...
	Id                int64     `json:"id" gorm:"PRIMARY_KEY";"AUTO_INCREMENT" ` // generated by DB
	CompanyId         int64     // belongs_to Company
	JobTitle          string    `gorm:"index:job_titles_idx" sql:"type:text;" json:"job_title"` // only required field
	PostingUrl        string    `sql:"type:text;" json:"posting_url"`
	Source            string    `sql:"type:text;" json:"source"`
	ReferredBy        string    `sql:"type:text;" json:"referred_by"`
	SalaryRange       string    `sql:"type:text;" json:"salary_range"`
	ApplicationMethod string    `sql:"type:text;" json:"application_method"`
	Note              string    `sql:"type:text;" json:"note"`  // markdown
	StartDate         time.Time `json:"start_date"`             // generated if not supplied
	CreatedAt         time.Time `json:"created_at"`             // generated if not supplied
	UpdatedAt         time.Time `json:"updated_at"`             // generated if not supplied
	People            []Person  `gorm:"many2many:jobs_people;"` // has and belongs to many people
	Company           Company   `gorm:"foreignkey:CompanyId"`
}

func (j Job) HolisticDeletion(db *gorm.DB) (bool, error) {

	transaction := db.Begin()
	// find all the events
	// returns slice of IEvent objects
	iEvents, _ := GetIEvents(db, j.Id)
	// delete all the events
	for _, event := range iEvents {
		err := db.Delete(event).Error
		if err != nil {
			transaction.Rollback()
			return false, err
		}
	}
	// delete the job
	if err := db.Delete(&j).Error; err != nil {
		transaction.Commit()
		return true, nil
	} else {
		transaction.Rollback()
		return false, err
	}
}

func (j Job) TransactionlessHolisticDeletion(db *gorm.DB) (bool, error) {
	iEvents, _ := GetIEvents(db, j.Id)
	// delete all the events
	for _, event := range iEvents {
		err := db.Delete(event).Error
		if err != nil {
			return false, err
			return false, err
		}
	}
	// delete the job
	if err := db.Delete(&j).Error; err != nil {
		return true, nil
	} else {
		return false, err
	}
}
