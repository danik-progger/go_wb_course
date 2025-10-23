package repos

import (
	"calendar/entities"
	"time"
)

type EventsRepo interface {
	GetEvent(user_id, event_id entities.Id) (entities.Event, error)
	GetEventsInRange(user_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error)
	AddEvent(user_id entities.Id, e entities.Event) (entities.Id, error)
	UpdateEvent(user_id, event_id entities.Id, e entities.Event) error
	DeleteEvent(user_id, id entities.Id) error
}
