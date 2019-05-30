package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// ┌──────────────────────────────────────────┐
// │Event                                     │
// │-------------------------------------     │
// │id:          int                          │
// │job_id    :  int                          │
// │created_at:  date_time                    │
// │people_ids:  array of ints (optional)     │
// │note:        string (optional - markdown) │
// │                                          │
// │belongs_to:  job                          │
// │has_and_belongs_to_many: people           │
// └──────────────────────────────────────────┘

type Event struct {
	Id        int64     `json:"id" gorm:"PRIMARY_KEY";"AUTO_INCREMENT" ` // generated by DB
	CreatedAt time.Time `json:"created_at"`                              // generated if not supplied
	UpdatedAt time.Time `json:"updated_at"`                              // generated if not supplied

	Note  string `sql:"type:text;" json:"note"` // markdown
	JobId int64  `json:"job_id"`
	Job   Job    `gorm:"foreignkey:JobId" json:"-"`
}

// implementing IEvent interface
func (e Event) CreationDate() time.Time {
	return e.CreatedAt
}

func (e Event) HolisticDeletion(db *gorm.DB) (bool, error) {

	transaction := db.Begin()
	if err := db.Delete(&e).Error; err != nil {
		transaction.Commit()
		return true, nil
	} else {
		transaction.Rollback()
		return false, err
	}

}
