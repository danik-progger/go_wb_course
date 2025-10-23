package app

import (
	repos "calendar/app/repos"
	"calendar/entities"
	"fmt"
	"time"
)

type Calendar[UserDB repos.UsersRepo, EventsDB repos.EventsRepo] struct {
	Users  UserDB
	Events EventsDB
}

func (c *Calendar[UserDB, EventsDB]) GetUser(id entities.Id) (entities.User, error) {
	return c.Users.GetUser(id)
}

func (c *Calendar[UserDB, EventsDB]) AddEvent(u entities.User, e entities.Event) (entities.Id, error) {
	exists, err := c.Users.HasUser(u.Id)
	if err != nil {
		return -1, fmt.Errorf("Failed to find out if user exists")
	}

	if !exists {
		if err := c.Events.AddUser(u.Id); err != nil {
			return -1, fmt.Errorf("Failed to add user")
		}
		if err := c.Users.AddUser(u); err != nil {
			return -1, fmt.Errorf("Failed to add user")
		}
	}

	return c.Events.AddEvent(u.Id, e)
}

func (c *Calendar[UserDB, EventsDB]) UpdateEvent(u entities.User, e entities.Event) error {
	if exists, err := c.Users.HasUser(u.Id); !exists || err != nil {
		return fmt.Errorf("Failed to find user")
	}

	return c.Events.UpdateEvent(u.Id, e.Id, e)
}

func (c *Calendar[UserDB, EventsDB]) DeleteEvent(u entities.User, e entities.Event) error {
	if exists, err := c.Users.HasUser(u.Id); !exists || err != nil {
		return fmt.Errorf("Failed to find user")
	}

	return c.Events.DeleteEvent(u.Id, e.Id)
}

func (c *Calendar[UserDB, EventsDB]) GetEvent(u entities.User, e_id entities.Id) (entities.Event, error) {
	return c.Events.GetEvent(u.Id, e_id)
}

func (c *Calendar[UserDB, EventsDB]) EventsInRange(u entities.User, from time.Time, to time.Time) ([]entities.Event, error) {
	if exists, err := c.Users.HasUser(u.Id); !exists || err != nil {
		return nil, fmt.Errorf("Failed to find user")
	}

	return c.Events.GetEventsInRange(u.Id, from, to)
}
