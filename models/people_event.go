package models

import "encoding/json"

type PeopleEvent struct {
	Event
	People []Person
}

func (pe *PeopleEvent) MarshalJSON() ([]byte, error) {
	person_ids := make([]int64, len(pe.People))
	for i, p := range pe.People {
		person_ids[i] = p.Id
	}
	type Alias PeopleEvent
	return json.Marshal(&struct {
		PersonIDs []int64 `json:"person_ids"`
		*Alias
	}{
		Alias:     (*Alias)(pe),
		PersonIDs: person_ids,
	})
}
