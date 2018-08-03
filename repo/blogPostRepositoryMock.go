package repo

import (
	"github.com/Ullaakut/Bloggo/model"
	"github.com/stretchr/testify/mock"
)

// BlogPostRepositoryMock is a mock of BlogPostRepository
type BlogPostRepositoryMock struct {
	mock.Mock
}

// Store mock
func (m *BlogPostRepositoryMock) Store(content *model.BlogPost) (*model.BlogPost, error) {
	args := m.Called(content)
	return args.Get(0).(*model.BlogPost), args.Error(1)
}

// Retrieve mock
func (m *BlogPostRepositoryMock) Retrieve(id uint) (model.BlogPost, error) {
	args := m.Called(id)
	return args.Get(0).(model.BlogPost), args.Error(1)
}

// TODO: Add Update
// TODO: Add Delete
// TODO: Add RetrieveAll
// TODO: Add Find? Retrieve with filters could be cool (filter by id, author, etc.)
