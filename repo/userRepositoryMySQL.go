package repo

import (
	"github.com/Ullaakut/Bloggo/model"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

// UserRepositoryMySQL is a repository to manage users stored using Gorm
type UserRepositoryMySQL struct {
	db *gorm.DB

	log *zerolog.Logger
}

// NewUserRepositoryMySQL creates a new user repository using the given gorm DB as backend
func NewUserRepositoryMySQL(log *zerolog.Logger, db *gorm.DB) *UserRepositoryMySQL {
	return &UserRepositoryMySQL{
		db: db,

		log: log,
	}
}

// Retrieve returns the user with the given ID from the database
func (r *UserRepositoryMySQL) Retrieve(token string) (bool, error) {
	user := model.User{
		Token: token,
	}

	err := r.db.Where(&user).First(&user).Error

	return user.Admin, err
}
