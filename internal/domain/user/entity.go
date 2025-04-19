package user

type User struct {
	ID   string
	Name string
}

func NewUser(id string, name string) *User {
	return &User{ID: id, Name: name}
}
