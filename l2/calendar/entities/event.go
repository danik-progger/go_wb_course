package entities

import "time"

type Event struct {
	Id          Id        `json:"id"`
	StartingAt  time.Time `json:"starting_at"`
	EndingAt    time.Time `json:"ending_at"`
	Description string    `json:"description"`
}

func NewEvent(id Id, dateStr string, desc string) (Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Event{}, err
	}
	return Event{
		Id:          id,
		StartingAt:  date,
		EndingAt:    date,
		Description: desc,
	}, nil
}
