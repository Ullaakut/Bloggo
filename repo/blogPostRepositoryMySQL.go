package repo

import (
	errortypes "github.com/Ullaakut/Bloggo/errorTypes"
	"github.com/Ullaakut/Bloggo/model"

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
func (r *BlogPostRepositoryMySQL) Store(content *model.BlogPost) (*model.BlogPost, error) {
	err := r.db.Create(content).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// if the error is of type duplicate entry
		if mysqlError.Number == 1062 {
			return nil, errortypes.ErrDuplicateEntry
		}
	}

	return content, err
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

// TODO: Add Update
// TODO: Add Delete
// TODO: Add Find? Retrieve with filters could be cool (filter by id, author, etc.)
