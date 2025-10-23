package server

import (
	"calendar/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	u_id := s.router.GetUrlParam(r, "user_id")
	user_id, err := strconv.ParseInt(u_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user_id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id))
	if err != nil {
		fmt.Printf("🔴 Failed to find user with id: %d\n", user_id)
		return
	}

	var event entities.Event
	json.NewDecoder(r.Body).Decode(&event)
	s.cal.AddEvent(u, event)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	u_id := s.router.GetUrlParam(r, "user_id")
	user_id, err := strconv.ParseInt(u_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user_id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id))
	if err != nil {
		fmt.Printf("🔴 Failed to find user with id: %d\n", user_id)
		return
	}

	var event entities.Event
	json.NewDecoder(r.Body).Decode(&event)
	s.cal.UpdateEvent(u, event)
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	u_id := s.router.GetUrlParam(r, "user_id")
	user_id, err := strconv.ParseInt(u_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user_id")
		return
	}
	event_id := s.router.GetUrlParam(r, "event_id")
	e_id, err := strconv.ParseInt(event_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user_id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id))
	if err != nil {
		fmt.Printf("🔴 Failed to find user with id: %d\n", user_id)
		return
	}

	e, err := s.cal.GetEvent(u, e_id)

	s.cal.DeleteEvent(u, e)
}

func (s *Server) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	starting_at := r.URL.Query().Get("date")
	starting_at_date, err := time.Parse("2006-01-02", starting_at)
	if err != nil {
		fmt.Println("🔴 Failed to parse date")
		return
	}
	ending_at := starting_at_date.AddDate(0, 0, 1)

	user_id := s.router.GetUrlParam(r, "user_id")
	user_id_int, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id_int))
	if err != nil {
		fmt.Println("🔴 Failed to find user by id")
		return
	}
	events, err := s.cal.EventsInRange(u, starting_at_date, ending_at)
	if err != nil {
		fmt.Println("🔴 Failed to find events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	starting_at := r.URL.Query().Get("date")
	starting_at_date, err := time.Parse("2006-01-02", starting_at)
	if err != nil {
		fmt.Println("🔴 Failed to parse date")
		return
	}
	ending_at := starting_at_date.AddDate(0, 0, 7)

	user_id := s.router.GetUrlParam(r, "user_id")
	user_id_int, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id_int))
	if err != nil {
		fmt.Println("🔴 Failed to find user by id")
		return
	}
	events, err := s.cal.EventsInRange(u, starting_at_date, ending_at)
	if err != nil {
		fmt.Println("🔴 Failed to find events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	starting_at := r.URL.Query().Get("date")
	starting_at_date, err := time.Parse("2006-01-02", starting_at)
	if err != nil {
		fmt.Println("🔴 Failed to parse date")
		return
	}
	ending_at := starting_at_date.AddDate(0, 1, 0)

	user_id := s.router.GetUrlParam(r, "user_id")
	user_id_int, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		fmt.Println("🔴 Failed to parse user id")
		return
	}
	u, err := s.cal.GetUser(entities.Id(user_id_int))
	if err != nil {
		fmt.Println("🔴 Failed to find user by id")
		return
	}
	events, err := s.cal.EventsInRange(u, starting_at_date, ending_at)
	if err != nil {
		fmt.Println("🔴 Failed to find events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
