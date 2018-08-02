package repo

import "github.com/stretchr/testify/mock"

// UserRepositoryMock is a mock of UserRepository
type UserRepositoryMock struct {
	mock.Mock
}

// Retrieve mock
func (m *UserRepositoryMock) Retrieve(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
