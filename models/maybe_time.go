package models

import (
	"time"
)

type MaybeTime struct {
	Just time.Time
	Err  error
}

func (m MaybeTime) Bind(f func(t time.Time) MaybeTime) MaybeTime {
	if m.IsError() {
		return m
	}
	return f(m.Just)
}
func (m MaybeTime) IsError() bool {
	return m.Err == nil
}

func JustTime(t time.Time) MaybeTime {
	return MaybeTime{Just: t}
}
func ErrorTime(e error) MaybeTime {
	return MaybeTime{Err: e}
}
