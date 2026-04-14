package server

import (
	"calendar/app"
	"calendar/app/repos"
	"calendar/entities"
	"calendar/infra/connections"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

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
	id := entities.Id(len(m.events) + 1)
	ev := e
	ev.Id = id
	m.events[id] = ev
	return id, nil
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
		if !event.StartingAt.Before(from) && event.StartingAt.Before(to) {
			result = append(result, event)
		}
	}
	return result, nil
}

func TestHandlers(t *testing.T) {
	setup := func() *httptest.Server {
		usersRepo := &MockUsersRepo{}
		eventsRepo := &MockEventsRepo{}
		cal := app.Calendar[repos.UsersRepo, repos.EventsRepo]{
			Users:  usersRepo,
			Events: eventsRepo,
		}

		r := chi.NewRouter()
		server := &Server{cal: cal, router: &connections.Router{}}
		r.Post("/create_event", server.CreateEvent)
		r.Post("/update_event", server.UpdateEvent)
		r.Post("/delete_event", server.DeleteEvent)
		r.Get("/events_for_day", server.GetEventsForDay)
		r.Get("/events_for_week", server.GetEventsForWeek)
		r.Get("/events_for_month", server.GetEventsForMonth)

		return httptest.NewServer(r)
	}

	t.Run("CreateEvent", func(t *testing.T) {
		ts := setup()
		defer ts.Close()

		form := url.Values{}
		form.Set("user_id", "1")
		form.Set("date", time.Now().Format("2006-01-02"))
		form.Set("event", "Test Event")

		resp, err := http.PostForm(ts.URL+"/create_event", form)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.StatusCode)
		}

		var result jsonResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
		if result.Result == "" {
			t.Error("expected non-empty result")
		}
	})

	t.Run("GetEventsForDay", func(t *testing.T) {
		usersRepo := &MockUsersRepo{users: map[entities.Id]entities.User{
			1: {Id: 1, Name: "test"},
		}}
		eventsRepo := &MockEventsRepo{}
		cal := app.Calendar[repos.UsersRepo, repos.EventsRepo]{
			Users:  usersRepo,
			Events: eventsRepo,
		}

		r := chi.NewRouter()
		server := &Server{cal: cal, router: &connections.Router{}}
		r.Get("/events_for_day", server.GetEventsForDay)

		ts := httptest.NewServer(r)
		defer ts.Close()

		today := time.Now().Format("2006-01-02")
		req, _ := http.NewRequest("GET", ts.URL+"/events_for_day?user_id=1&date="+today, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.StatusCode)
		}
	})

	t.Run("CreateEvent missing params returns 400", func(t *testing.T) {
		ts := setup()
		defer ts.Close()

		resp, err := http.PostForm(ts.URL+"/create_event", url.Values{})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status Bad Request, got %v", resp.StatusCode)
		}

		var result jsonResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
		if result.Error == "" {
			t.Error("expected error message")
		}
	})
}
