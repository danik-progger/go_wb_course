package entities

import "time"

type Event struct {
	id          Id
	starting_at time.Time
	ending_at   time.Time
	description string
}

func (e *Event) Id() Id {
	return e.id
}

func (e *Event) StartingAt() time.Time {
	return e.starting_at
}

func (e *Event) EndingAt() time.Time {
	return e.ending_at
}

func (e *Event) Desc() string {
	return e.description
}

func NewEvent(id Id, from, to time.Time, desc string) Event {
	return Event{
		id:          id,
		starting_at: from,
		ending_at:   to,
		description: desc,
	}
}
