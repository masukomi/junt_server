package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

func EmptyStringForNilString(value interface{}) string {
	if value == nil {
		return ""
	} else {
		return value.(string)
	}
}
func MaybeTimeFromValue(value string) MaybeTime {
	newTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return JustTime(newTime)
	} else {
		return ErrorTime(errors.New("invalid created_at time"))
	}
}

func MaybeCompanyFromId(companyId int64, db *gorm.DB) MaybeCompany {
	company := Company{}
	if err := db.First(&company, companyId).Error; err != nil {
		return JustCompany(company)
	}
	return ErrorCompany(errors.New("invalid company id"))
}

func ExtractIdsFromJsonArray(arr []interface{}) []int64 {
	theIds := []int64{}
	for _, num := range arr { // []interface{}
		theIds = append(theIds, int64(num.(float64)))
	}
	return theIds
}
