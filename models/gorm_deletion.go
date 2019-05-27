package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func DeleteOrRollback(transaction *gorm.DB, db *gorm.DB, x interface{}) (bool, error) {
	err := db.Delete(x).Error
	if err != nil {
		transaction.Rollback()
		return false, err
	}
	return true, nil

}
