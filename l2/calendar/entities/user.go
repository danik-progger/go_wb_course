package entities

type User struct {
	id   Id
	name string
}

func (u *User) Id() Id {
	return u.id
}

func NewUser(id Id, name string) User {
	return User{
		id:   id,
		name: name,
	}
}

func (u *User) GetName() string {
	return u.name
}
