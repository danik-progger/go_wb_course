package app

import (
	repos "calendar/app/repos"
	"calendar/entities"
	"errors"
	"testing"
	"time"
)

type MockUsersRepo struct {
	users map[entities.Id]entities.User
}

func (m *MockUsersRepo) AddUser(u entities.User) error {
	if m.users == nil {
		m.users = make(map[entities.Id]entities.User)
	}
	m.users[u.Id] = u
	return nil
}
func (m *MockUsersRepo) GetUser(id entities.Id) (entities.User, error) {
	user, ok := m.users[id]
	if !ok {
		return entities.User{}, errors.New("user not found")
	}
	return user, nil
}
func (m *MockUsersRepo) HasUser(id entities.Id) (bool, error) {
	_, ok := m.users[id]
	return ok, nil
}

type MockEventsRepo struct {
	events map[entities.Id][]entities.Event
}

func (m *MockEventsRepo) AddUser(u_id entities.Id) error {
	return nil
}

func (m *MockEventsRepo) AddEvent(u_id entities.Id, e entities.Event) (entities.Id, error) {
	if m.events == nil {
		m.events = make(map[entities.Id][]entities.Event)
	}
	id := entities.Id(len(m.events[u_id]) + 1)
	e.Id = id
	m.events[u_id] = append(m.events[u_id], e)
	return id, nil
}
func (m *MockEventsRepo) UpdateEvent(u_id, e_id entities.Id, e entities.Event) error {
	for i, ev := range m.events[u_id] {
		if ev.Id == e_id {
			m.events[u_id][i] = e
			return nil
		}
	}
	return nil
}
func (m *MockEventsRepo) DeleteEvent(u_id, e_id entities.Id) error {
	events := m.events[u_id]
	for i, ev := range events {
		if ev.Id == e_id {
			m.events[u_id] = append(events[:i], events[i+1:]...)
			return nil
		}
	}
	return nil
}
func (m *MockEventsRepo) GetEvent(u_id, e_id entities.Id) (entities.Event, error) {
	for _, e := range m.events[u_id] {
		if e.Id == e_id {
			return e, nil
		}
	}
	return entities.Event{}, errors.New("event not found")
}
func (m *MockEventsRepo) GetEventsInRange(u_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error) {
	var result []entities.Event
	for _, e := range m.events[u_id] {
		if !e.StartingAt.Before(from) && e.StartingAt.Before(to) {
			result = append(result, e)
		}
	}
	return result, nil
}

func TestCalendar_AddEvent_NewUser(t *testing.T) {
	usersRepo := &MockUsersRepo{}
	eventsRepo := &MockEventsRepo{}
	cal := Calendar[repos.UsersRepo, repos.EventsRepo]{Users: usersRepo, Events: eventsRepo}

	dateStr := time.Now().Format("2006-01-02")
	event, _ := entities.NewEvent(0, dateStr, "Test Event")
	user := entities.NewUser(1, "testuser")

	id, err := cal.AddEvent(user, event)

	if err != nil {
		t.Fatalf("AddEvent failed: %v", err)
	}
	if id != 1 {
		t.Errorf("expected event id 1, got %v", id)
	}
	hasUser, _ := usersRepo.HasUser(user.Id)
	if !hasUser {
		t.Error("expected user to be added")
	}
}
