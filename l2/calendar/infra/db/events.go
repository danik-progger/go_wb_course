package db

import (
	"calendar/entities"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type EventsDb struct{}

func InitEventsDb() (*EventsDb, error) {
	return &EventsDb{}, nil
}

func ConnectToUser(u_id entities.Id) (*sqlx.DB, error) {
	cal := fmt.Sprintf("./calendars/calendar_%d", u_id)
	return sqlx.Connect("sqlite", cal)
}

func (db *EventsDb) AddUser(u_id entities.Id) error {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return err
	}
	query := `CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		starting_at DATETIME,
		ending_at DATETIME,
		description TEXT
	)`
	_, err = cal.Exec(query)
	return err
}

func (db *EventsDb) GetEvent(u_id, e_id entities.Id) (entities.Event, error) {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return entities.Event{}, err
	}
	var event entities.Event
	query := "SELECT id, starting_at, ending_at, description FROM events WHERE id = ?"
	err = cal.Get(&event, query, e_id)
	return event, err
}

func (db *EventsDb) GetEventsInRange(u_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error) {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return nil, err
	}
	var events []entities.Event
	query := "SELECT id, starting_at, ending_at, description FROM events WHERE starting_at <= ? AND ending_at >= ?"
	err = cal.Select(&events, query, to, from)
	return events, err
}
func (db *EventsDb) AddEvent(u_id entities.Id, e entities.Event) (entities.Id, error) {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return -1, err
	}
	query := "INSERT INTO events (starting_at, ending_at, description) VALUES (?, ?, ?)"
	res, err := cal.Exec(query, e.StartingAt, e.EndingAt, e.Desc)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return entities.Id(id), err
}
func (db *EventsDb) UpdateEvent(u_id, e_id entities.Id, e entities.Event) error {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return err
	}
	query := "UPDATE events SET starting_at = ?, ending_at = ?, description = ? WHERE id = ?"
	_, err = cal.Exec(query, e.StartingAt, e.EndingAt, e.Desc, e_id)
	return err
}
func (db *EventsDb) DeleteEvent(u_id, e_id entities.Id) error {
	cal, err := ConnectToUser(u_id)
	if err != nil {
		return err
	}
	query := "DELETE FROM events WHERE id = ?"
	_, err = cal.Exec(query, e_id)
	return err
}
