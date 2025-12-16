package app

import (
	repos "calendar/app/repos"
	"calendar/entities"
	"errors"
	"testing"
	"time"
)

// Simplified MockUsersRepo for debugging.
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

// Simplified MockEventsRepo for debugging.
type MockEventsRepo struct{}

func (m *MockEventsRepo) AddUser(u_id entities.Id) error {
	return nil
}

func (m *MockEventsRepo) AddEvent(u_id entities.Id, e entities.Event) (entities.Id, error) {
	return 1, nil
}
func (m *MockEventsRepo) UpdateEvent(u_id, e_id entities.Id, e entities.Event) error { return nil }
func (m *MockEventsRepo) DeleteEvent(u_id, e_id entities.Id) error                   { return nil }
func (m *MockEventsRepo) GetEvent(u_id, e_id entities.Id) (entities.Event, error) {
	return entities.Event{}, nil
}
func (m *MockEventsRepo) GetEventsInRange(u_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error) {
	return nil, nil
}

func TestCalendar_AddEvent_NewUser(t *testing.T) {
	usersRepo := &MockUsersRepo{}
	eventsRepo := &MockEventsRepo{}
	cal := Calendar[repos.UsersRepo, repos.EventsRepo]{Users: usersRepo, Events: eventsRepo}
	user := entities.NewUser(1, "testuser")
	event := entities.NewEvent(0, time.Time{}, time.Time{}, "Test Event")

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
