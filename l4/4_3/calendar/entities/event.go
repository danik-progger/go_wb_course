package entities

import "time"

type Event struct {
	Id          Id        `db:"id"`
	StartingAt  time.Time `db:"starting_at"`
	EndingAt    time.Time `db:"ending_at"`
	Description string    `db:"description"`
}


func NewEvent(id Id, from, to time.Time, desc string) Event {
	return Event{
		Id:          id,
		StartingAt:  from,
		EndingAt:    to,
		Description: desc,
	}
}
