package server

import (
	"calendar/app"
	"calendar/app/repos"
	"calendar/infra/connections"
	"calendar/infra/db"
	logger "calendar/infra/log"
	"fmt"
)

type Server struct {
	cal    app.Calendar[repos.UsersRepo, repos.EventsRepo]
	router *connections.Router
}

func NewServer(port string) *Server {
	r := connections.NewRouter(port)
	usersDB := db.InitUsersDb()
	eventsDB, err := db.InitEventsDb()
	if err != nil {
		fmt.Println("🔴 Failed to initialize events database")
		return nil
	}
	cal := app.Calendar[repos.UsersRepo, repos.EventsRepo]{
		Users:  usersDB,
		Events: eventsDB,
	}
	s := &Server{
		cal:    cal,
		router: r,
	}

	s.setUpRoutes()
	s.router.AddMiddleware(logger.LoggingMiddleware)

	return s
}

func (s *Server) Run() error {
	return s.router.Run()
}

func (s *Server) setUpRoutes() {
	s.router.Post("/users/{user_id}/events", s.CreateEvent)
	s.router.Post("/users/{user_id}/events/{event_id}", s.UpdateEvent)
	s.router.Post("/users/{user_id}/events/{event_id}/delete", s.DeleteEvent) // Or use DELETE method
	s.router.Get("/users/{user_id}/events/day", s.GetEventsForDay)
	s.router.Get("/users/{user_id}/events/week", s.GetEventsForWeek)
	s.router.Get("/users/{user_id}/events/month", s.GetEventsForMonth)
}
