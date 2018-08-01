package repo

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// UserRepository is a repository to manage users stored using Gorm
type UserRepository struct {
	db *gorm.DB

	// TODO: Remove and use real DB
	tmp map[string]struct{}
}

// NewUserRepository creates a new user repository using the given gorm DB as backend
func NewUserRepository() *UserRepository { //db *gorm.DB) *UserRepository {
	return &UserRepository{
		// db: db,

		// TODO: Remove and use real DB
		tmp: map[string]struct{}{
			"auth0|596f27c2c3709661e9cea37d": struct{}{},
		},
	}
}

// Retrieve returns the user with the given ID from the database
func (r *UserRepository) Retrieve(id string) error {
	_, ok := r.tmp[id]
	if !ok {
		return errors.New("user not found")
	}

	return nil

	// return r.db.First(nil, id).Error
}
