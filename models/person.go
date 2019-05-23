package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// ┌──────────────────────────────────────────┐
// │Person                                    │
// │----------------                          │
// │id:         int                           │
// │name:       string                        │
// │email:      string (optional)             │
// │phone:      string (optional)             │
// │created_at: date_time (optional)          │
// │updated_at: date_time (optional)          │
// │note:       string (optional - markdown)  │
// │                                          │
// │has_and_belongs_to_many:   events (maybe) │
// │has_and_belongs_to_many:   jobs (probably)│
// │                                          │
// │(events gets you interviews)              │
// └──────────────────────────────────────────┘
// dates are ISO 8601 RFC3339
// t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

type Person struct {
	gorm.Model
	// Id, CreatedAt, UpdatedAt would normally come from gorm.Model but we need to specify
	// the json keys for them so...
	Id           int64         `json:"id" gorm:"PRIMARY_KEY";"AUTO_INCREMENT" `             // generated by DB
	Name         string        `gorm:"index:people_names_idx" sql:"type:text;" json:"name"` // only required field
	Email        string        `sql:"type:text;" json:"email"`
	Phone        string        `sql:"type:text;" json:"phone"`
	Note         string        `sql:"type:text;" json:"note"` // markdown
	CompanyId    int64         // belongs_to Company
	CreatedAt    time.Time     `json:"created_at"`             // generated if not supplied
	UpdatedAt    time.Time     `json:"updated_at"`             // generated if not supplied
	Jobs         []Job         `gorm:"many2many:jobs_people;"` // has and belongs to many jobs
	Company      Company       `gorm:"gorm:foreignkey:CompanyId"`
	ThanksEmails []ThanksEmail `gorm:"many2many:people_thanks_emails;"` // has and belongs to many jobs
}