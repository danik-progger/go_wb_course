package db

import "calendar/entities"

type UsersDb struct {
	storage map[entities.Id]entities.User
}

func InitUsersDb() *UsersDb {
	return &UsersDb{storage: make(map[entities.Id]entities.User)}
}

func (db *UsersDb) AddUser(u entities.User) error {
	db.storage[u.Id()] = u
	return nil
}

func (db *UsersDb) GetUser(id entities.Id) (entities.User, error) {
	return db.storage[id], nil
}

func (db *UsersDb) HasUser(id entities.Id) (bool, error) {
	_, ok := db.storage[id]
	return ok, nil
}
