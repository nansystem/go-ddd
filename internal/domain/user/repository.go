package user

type Repository interface {
	GetUsers() ([]*User, error)
	GetUserByID(id string) (*User, error)
	CreateUser(user *User) error
}
