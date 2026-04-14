package db

import (
	"calendar/entities"
	"sync"
	"time"
)

type EventsDb struct {
	mu     sync.RWMutex
	events map[entities.Id][]entities.Event
	nextID entities.Id
}

func InitEventsDb() (*EventsDb, error) {
	return &EventsDb{
		events: make(map[entities.Id][]entities.Event),
		nextID: 1,
	}, nil
}

func (db *EventsDb) AddUser(u_id entities.Id) error {
	return nil
}

func (db *EventsDb) GetEvent(u_id, e_id entities.Id) (entities.Event, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, e := range db.events[u_id] {
		if e.Id == e_id {
			return e, nil
		}
	}
	return entities.Event{}, nil
}

func (db *EventsDb) GetEventsInRange(u_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var result []entities.Event
	for _, e := range db.events[u_id] {
		if !e.StartingAt.Before(from) && e.StartingAt.Before(to) {
			result = append(result, e)
		}
	}
	return result, nil
}

func (db *EventsDb) AddEvent(u_id entities.Id, e entities.Event) (entities.Id, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := db.nextID
	db.nextID++
	e.Id = id
	db.events[u_id] = append(db.events[u_id], e)
	return id, nil
}

func (db *EventsDb) UpdateEvent(u_id, e_id entities.Id, e entities.Event) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, ev := range db.events[u_id] {
		if ev.Id == e_id {
			db.events[u_id][i] = e
			return nil
		}
	}
	return nil
}

func (db *EventsDb) DeleteEvent(u_id, e_id entities.Id) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	events := db.events[u_id]
	for i, ev := range events {
		if ev.Id == e_id {
			db.events[u_id] = append(events[:i], events[i+1:]...)
			return nil
		}
	}
	return nil
}
