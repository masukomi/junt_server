package models

type Followup struct {
	Event
	People []Person `gorm:"many2many:followups_people;"`
}

// proof it implements the interface
// var _ IEvent = (*Followup)(nil)
