package usecase

import (
	"github.com/stretchr/testify/mock"

	"github.com/nansystem/go-ddd/internal/domain/user"
)

// MockUserService はUserServiceのモック実装です
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers() ([]*user.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) CreateUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}
