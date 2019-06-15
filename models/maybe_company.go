package models

type MaybeCompany struct {
	Just Company
	Err  error
}

func (m MaybeCompany) Bind(f func(c Company) MaybeCompany) MaybeCompany {
	if m.IsError() {
		return m
	}
	return f(m.Just)
}
func (m MaybeCompany) IsError() bool {
	return m.Err == nil
}

func JustCompany(c Company) MaybeCompany {
	return MaybeCompany{Just: c}
}
func ErrorCompany(e error) MaybeCompany {
	return MaybeCompany{Err: e}
}
