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
	s.router.Post("/create_event", s.CreateEvent)
	s.router.Post("/update_event", s.UpdateEvent)
	s.router.Post("/delete_event", s.DeleteEvent)
	s.router.Get("/events_for_day", s.GetEventsForDay)
	s.router.Get("/events_for_week", s.GetEventsForWeek)
	s.router.Get("/events_for_month", s.GetEventsForMonth)
}
