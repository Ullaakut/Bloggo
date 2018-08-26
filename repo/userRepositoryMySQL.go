package repo

import (
	"github.com/Ullaakut/Bloggo/errortype"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/go-sql-driver/mysql"
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
func (r *UserRepositoryMySQL) Retrieve(user *model.User) (*model.User, error) {
	err := r.db.Where(user).First(user).Error
	return user, err
}

// AdminExists returns true if an admin exists, false otherwise
func (r *UserRepositoryMySQL) AdminExists() bool {
	filter := &model.User{
		IsAdmin: true,
	}

	return r.db.Where(filter).First(filter).Error == nil
}

// Store saves a new user in the database.
func (r *UserRepositoryMySQL) Store(user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// if the error is of type duplicate entry
		if mysqlError.Number == 1062 {
			return nil, errortype.ErrDuplicateEntry
		}
	}

	return user, err
}
