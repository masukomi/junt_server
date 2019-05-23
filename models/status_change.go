package models

type StatusChange struct {
	Event
	From string `sql:"type:text" json:"from"`
	To   string `sql:"type:text" json:"to"`
}
