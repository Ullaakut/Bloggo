package repo

import "github.com/stretchr/testify/mock"

// UserRepositoryMock is a mock of UserRepository
type UserRepositoryMock struct {
	mock.Mock
}

// Retrieve mock
func (m *UserRepositoryMock) Retrieve(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
