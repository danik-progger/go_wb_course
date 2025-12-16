package server

import (
	"bytes"
	"calendar/app"
	"calendar/app/repos"
	"calendar/entities"
	"calendar/infra/connections"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

// MockUsersRepo is a mock of UsersRepo for testing.
type MockUsersRepo struct {
	users map[entities.Id]entities.User
	err   error
}

func (m *MockUsersRepo) AddUser(u entities.User) error {
	if m.err != nil {
		return m.err
	}
	if m.users == nil {
		m.users = make(map[entities.Id]entities.User)
	}
	m.users[u.Id] = u
	return nil
}

func (m *MockUsersRepo) GetUser(id entities.Id) (entities.User, error) {
	if m.err != nil {
		return entities.User{}, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return entities.User{}, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUsersRepo) HasUser(id entities.Id) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	_, ok := m.users[id]
	return ok, nil
}

// MockEventsRepo is a mock of EventsRepo for testing.
type MockEventsRepo struct {
	events map[entities.Id]entities.Event
	err    error
}

func (m *MockEventsRepo) AddUser(u_id entities.Id) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *MockEventsRepo) AddEvent(u_id entities.Id, e entities.Event) (entities.Id, error) {
	if m.err != nil {
		return -1, m.err
	}
	if m.events == nil {
		m.events = make(map[entities.Id]entities.Event)
	}
	newEvent := entities.NewEvent(entities.Id(len(m.events)+1), e.StartingAt, e.EndingAt, e.Description)
	m.events[newEvent.Id] = newEvent
	return newEvent.Id, nil
}

func (m *MockEventsRepo) UpdateEvent(u_id, e_id entities.Id, e entities.Event) error {
	if m.err != nil {
		return m.err
	}
	m.events[e_id] = e
	return nil
}

func (m *MockEventsRepo) DeleteEvent(u_id, e_id entities.Id) error {
	if m.err != nil {
		return m.err
	}
	delete(m.events, e_id)
	return nil
}

func (m *MockEventsRepo) GetEvent(u_id, e_id entities.Id) (entities.Event, error) {
	if m.err != nil {
		return entities.Event{}, m.err
	}
	event, ok := m.events[e_id]
	if !ok {
		return entities.Event{}, errors.New("event not found")
	}
	return event, nil
}

func (m *MockEventsRepo) GetEventsInRange(u_id entities.Id, from time.Time, to time.Time) ([]entities.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []entities.Event
	for _, event := range m.events {
		if !event.StartingAt.After(to) && !event.EndingAt.Before(from) {
			result = append(result, event)
		}
	}
	return result, nil
}

func TestHandlers(t *testing.T) {
	setup := func() (*httptest.Server, *MockUsersRepo, *MockEventsRepo) {
		usersRepo := &MockUsersRepo{}
		eventsRepo := &MockEventsRepo{}
		cal := app.Calendar[repos.UsersRepo, repos.EventsRepo]{
			Users:  usersRepo,
			Events: eventsRepo,
		}

		r := chi.NewRouter()
		server := &Server{cal: cal, router: &connections.Router{}}
		r.Post("/users/{user_id}/events", server.CreateEvent)
		r.Post("/users/{user_id}/events/{event_id}", server.UpdateEvent)
		r.Post("/users/{user_id}/events/{event_id}/delete", server.DeleteEvent)
		r.Get("/users/{user_id}/events/day", server.GetEventsForDay)
		r.Get("/users/{user_id}/events/week", server.GetEventsForWeek)
		r.Get("/users/{user_id}/events/month", server.GetEventsForMonth)

		ts := httptest.NewServer(r)

		user := entities.NewUser(1, "testuser")
		usersRepo.AddUser(user)

		return ts, usersRepo, eventsRepo
	}

	t.Run("CreateEvent", func(t *testing.T) {
		ts, _, _ := setup()
		defer ts.Close()

		event := entities.NewEvent(0, time.Now(), time.Now().Add(time.Hour), "Test Event")
		body, _ := json.Marshal(event)

		req, _ := http.NewRequest("POST", ts.URL+"/users/1/events", bytes.NewReader(body))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.StatusCode)
		}
	})

	t.Run("GetEventsForDay", func(t *testing.T) {
		ts, usersRepo, eventsRepo := setup()
		defer ts.Close()

		// Add an event for the test
		eventsRepo.AddEvent(usersRepo.users[1].Id, entities.NewEvent(1, time.Now(), time.Now().Add(time.Hour), "Test Event"))

		req, _ := http.NewRequest("GET", ts.URL+"/users/1/events/day?date="+time.Now().Format("2006-01-02"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.StatusCode)
		}

		var events []entities.Event
		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
			t.Fatal(err)
		}

		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}
	})
}
