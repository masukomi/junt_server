package models

type MaybeJob struct {
	Just Job
	Err  error
}

func (m MaybeJob) Bind(f func(j Job) MaybeJob) MaybeJob {
	if m.IsError() {
		return m
	}
	return f(m.Just)
}
func (m MaybeJob) IsError() bool {
	return m.Err == nil
}

func JustJob(j Job) MaybeJob {
	return MaybeJob{Just: j}
}
func ErrorJob(e error) MaybeJob {
	return MaybeJob{Err: e}
}
