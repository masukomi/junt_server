package models

import (
	"time"
)

type Homework struct {
	Event
	DueDate time.Time `json:"due_date"` // generated if not supplied
}
