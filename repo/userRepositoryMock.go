package repo

import (
	"github.com/Ullaakut/Bloggo/model"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock of UserRepository
type UserRepositoryMock struct {
	mock.Mock
}

// Retrieve mock
func (m *UserRepositoryMock) Retrieve(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Store mock
func (m *UserRepositoryMock) Store(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// AdminExists mock
func (m *UserRepositoryMock) AdminExists() bool {
	args := m.Called()
	return args.Bool(0)
}
