package server

import (
	"calendar/entities"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type jsonResponse struct {
	Result string           `json:"result,omitempty"`
	Error  string           `json:"error,omitempty"`
	Events []entities.Event `json:"events,omitempty"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(jsonResponse{Error: msg})
}

func writeResult(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse{Result: msg})
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		userIDStr = r.URL.Query().Get("user_id")
	}
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "missing user_id")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	dateStr := r.FormValue("date")
	if dateStr == "" {
		dateStr = r.URL.Query().Get("date")
	}
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "missing date")
		return
	}

	title := r.FormValue("event")
	if title == "" {
		title = r.URL.Query().Get("event")
	}

	event, err := entities.NewEvent(0, dateStr, title)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	u := entities.NewUser(entities.Id(userID), "")
	id, err := s.cal.AddEvent(u, event)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	writeResult(w, "event "+strconv.FormatInt(int64(id), 10)+" created")
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		userIDStr = r.URL.Query().Get("user_id")
	}
	eventIDStr := r.FormValue("event_id")
	if eventIDStr == "" {
		eventIDStr = s.router.GetUrlParam(r, "event_id")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid event_id")
		return
	}

	dateStr := r.FormValue("date")
	if dateStr == "" {
		dateStr = r.URL.Query().Get("date")
	}
	title := r.FormValue("event")
	if title == "" {
		title = r.URL.Query().Get("event")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	u := entities.NewUser(entities.Id(userID), "")
	event := entities.Event{
		Id:          entities.Id(eventID),
		StartingAt:  date,
		EndingAt:    date,
		Description: title,
	}

	if err := s.cal.UpdateEvent(u, event); err != nil {
		writeError(w, http.StatusServiceUnavailable, "failed to update event")
		return
	}

	writeResult(w, "event updated")
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		userIDStr = r.URL.Query().Get("user_id")
	}
	eventIDStr := r.FormValue("event_id")
	if eventIDStr == "" {
		eventIDStr = s.router.GetUrlParam(r, "event_id")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid event_id")
		return
	}

	u := entities.NewUser(entities.Id(userID), "")

	e, err := s.cal.GetEvent(u, entities.Id(eventID))
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "event not found")
		return
	}

	if err := s.cal.DeleteEvent(u, e); err != nil {
		writeError(w, http.StatusServiceUnavailable, "failed to delete event")
		return
	}

	writeResult(w, "event deleted")
}

func (s *Server) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.getEventsInRange(w, r, 1)
}

func (s *Server) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.getEventsInRange(w, r, 7)
}

func (s *Server) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.getEventsInRange(w, r, 30)
}

func (s *Server) getEventsInRange(w http.ResponseWriter, r *http.Request, days int) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "missing date")
		return
	}

	startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}
	endDate := startDate.AddDate(0, 0, days)

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	u, err := s.cal.GetUser(entities.Id(userID))
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "user not found")
		return
	}

	events, err := s.cal.EventsInRange(u, startDate, endDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
