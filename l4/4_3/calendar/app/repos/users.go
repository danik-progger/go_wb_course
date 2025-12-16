package repos

import "calendar/entities"

type UsersRepo interface {
	AddUser(u entities.User) error
	GetUser(id entities.Id) (entities.User, error)
	HasUser(id entities.Id) (bool, error)
}
