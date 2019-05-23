package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Offer struct {
	gorm.Model
	Id        int64     `json:"id" gorm:"PRIMARY_KEY";"AUTO_INCREMENT" ` // generated by DB
	Note      string    `sql:"type:text;" json:"note"`                   // markdown
	Status    string    `sql:"type:text" json:"status"`
	CreatedAt time.Time `json:"created_at"` // generated if not supplied
	UpdatedAt time.Time `json:"updated_at"` // generated if not supplied

	JobId int64 `json:"job_id"`
	Job   Job   `gorm:"foreignkey:JobId"`
}
