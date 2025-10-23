package entities

type User struct {
	Id   Id     `db:"id"`
	Name string `db:"name"`
}

func NewUser(id Id, name string) User {
	return User{
		Id:   id,
		Name: name,
	}
}

func (u *User) GetName() string {
	return u.Name
}
