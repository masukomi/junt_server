package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type IEvent interface {
	CreationDate() time.Time
	HolisticDeletion(db *gorm.DB) (bool, error)
}

// the folowing BS is because Go hates you
// and doesn't want you to have anything nice
// like Ruby's Comparable functionality.
// I'm sorry. Really.
type ByCreationDate []IEvent

func (e ByCreationDate) Len() int {
	return len(e)
}

func (e ByCreationDate) Less(i, j int) bool {
	return e[i].CreationDate().Unix() < e[j].CreationDate().Unix()
}

func (e ByCreationDate) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

/// end madness... for now
