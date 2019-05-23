package models

type Followup struct {
	Event
	People []Person `gorm:"many2many:followups_people;"`
}
