package models

type Followup struct {
	PeopleEvent
	People []Person `gorm:"many2many:followups_people;"`
}
