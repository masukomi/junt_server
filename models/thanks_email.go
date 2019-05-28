package models

type ThanksEmail struct {
	PeopleEvent
	People []Person `gorm:"many2many:people_thanks_emails;"` // has and belongs to many jobs
}
