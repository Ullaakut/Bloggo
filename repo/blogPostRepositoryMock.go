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
func (m *BlogPostRepositoryMock) Retrieve(id uint) (*model.BlogPost, error) {
	args := m.Called(id)
	return args.Get(0).(*model.BlogPost), args.Error(1)
}

// Find mock
func (m *BlogPostRepositoryMock) Find(contains *string, limit *uint) ([]*model.BlogPost, error) {
	args := m.Called(contains, limit)

	if args.Get(0).([]*model.BlogPost) != nil {
		return args.Get(0).([]*model.BlogPost), args.Error(1)
	}
	return nil, args.Error(1)
}

// Update mock
func (m *BlogPostRepositoryMock) Update(content *model.BlogPost) error {
	args := m.Called(content)
	return args.Error(0)
}

// Delete mock
func (m *BlogPostRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TODO: Add Find? Retrieve with filters could be cool (filter by id, author, etc.)
