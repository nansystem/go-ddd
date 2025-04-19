package usecase

import (
	"github.com/nansystem/go-ddd/internal/domain/user"
)

type UserService struct {
	userRepository user.Repository
}

func NewUserService(userRepository user.Repository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUsers() ([]*user.User, error) {
	return s.userRepository.GetUsers()
}

func (s *UserService) GetUserByID(id string) (*user.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) CreateUser(user *user.User) error {
	return s.userRepository.CreateUser(user)
}
