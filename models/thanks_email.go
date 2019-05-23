package models

type ThanksEmail struct {
	Event
	People []Person `gorm:"many2many:people_thanks_emails;"` // has and belongs to many jobs
}
