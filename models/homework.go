package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Homework struct {
	Event
	DueDate time.Time `json:"due_date"` // generated if not supplied
}

func (h *Homework) UpdateFromJson(data map[string]interface{}, db *gorm.DB) error {

	err := h.UpdateEventFromJson(data, db)
	if err != nil {
		if value, ok := data["due_date"]; ok {

			mTime := MaybeTimeFromValue(value.(string))
			if !mTime.IsError() {
				h.DueDate = mTime.Just
			} else {
				return errors.New("invalid created_at time: \"" + value.(string) + "\" Use RFC3339")
			}
		} else {
			return errors.New("Homework must have a valid due_date")
		}

	}
	return nil
}
