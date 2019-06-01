package models

import (
	"github.com/jinzhu/gorm"
)

type IJsonUpdateable interface {
	UpdateFromJson(data map[string]interface{}, db *gorm.DB) error
}
