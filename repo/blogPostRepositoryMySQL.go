package repo

import (
	errortypes "github.com/Ullaakut/Bloggo/errorTypes"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/pkg/errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

// BlogPostRepositoryMySQL is a repository to manage blog posts stored using Gorm
type BlogPostRepositoryMySQL struct {
	db *gorm.DB

	log *zerolog.Logger
}

// NewBlogPostRepositoryMySQL creates a new blog post repository using the given gorm DB as backend
func NewBlogPostRepositoryMySQL(log *zerolog.Logger, db *gorm.DB) *BlogPostRepositoryMySQL {
	return &BlogPostRepositoryMySQL{
		db: db,

		log: log,
	}
}

// Store saves a new blog post in the database.
func (r *BlogPostRepositoryMySQL) Store(post *model.BlogPost) (*model.BlogPost, error) {
	err := r.db.Create(post).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// if the error is of type duplicate entry
		if mysqlError.Number == 1062 {
			return nil, errortypes.ErrDuplicateEntry
		}
	}

	return post, err
}

// Retrieve returns the blog post with the given ID from the database
func (r *BlogPostRepositoryMySQL) Retrieve(id uint) (*model.BlogPost, error) {
	post := model.BlogPost{
		ID: id,
	}

	err := r.db.Where(&post).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return &post, errortypes.ErrNotFound
	}

	return &post, err
}

// RetrieveAll returns the blog post with the given ID from the database
func (r *BlogPostRepositoryMySQL) RetrieveAll() ([]*model.BlogPost, error) {
	var posts []*model.BlogPost

	err := r.db.Find(&posts).Error
	return posts, err
}

// Update saves a new blog post in the database.
func (r *BlogPostRepositoryMySQL) Update(post *model.BlogPost) error {
	var existingPost model.BlogPost

	// First, check if the blog post already exists
	// If not, return an error
	err := r.db.First(&existingPost, post.ID).Error
	if err == gorm.ErrRecordNotFound {
		return errortypes.ErrNotFound
	}
	if err != nil {
		return errors.Wrap(err, "could not get blog post from db")
	}

	post.CreatedAt = existingPost.CreatedAt

	// If it exists, saving this will overwrite the previous post
	// with the same ID
	err = r.db.Save(&post).Error
	if err != nil {
		return errors.Wrap(err, "could not save blog post in DB")
	}
	return nil
}

// TODO: Add Delete
// TODO: Add Find? Retrieve with filters could be cool (filter by id, author, etc.)
