package repo

import (
	"github.com/Ullaakut/Bloggo/model"
	"github.com/jinzhu/gorm"
)

// BlogPostRepository is a repository to manage blog posts stored using Gorm
type BlogPostRepository struct {
	db *gorm.DB

	// TODO: Remove and use real DB
	tmp        map[uint]model.BlogPost
	tmpIDCount uint
}

// NewBlogPostRepository creates a new blog post repository using the given gorm DB as backend
func NewBlogPostRepository() *BlogPostRepository { //db *gorm.DB) *BlogPostRepository {
	return &BlogPostRepository{
		// db: db,

		// TODO: Remove and use real DB
		tmp:        make(map[uint]model.BlogPost),
		tmpIDCount: 0,
	}
}

// Store saves a new blog post in the database.
func (r *BlogPostRepository) Store(content model.BlogPost) error {
	// return r.db.Create(content).Error
	r.tmp[r.tmpIDCount] = content
	r.tmpIDCount++
	return nil
}

// Retrieve returns the blog post with the given ID from the database
func (r *BlogPostRepository) Retrieve(id uint) (model.BlogPost, error) {
	content, ok := r.tmp[id]
	if !ok {
		return model.BlogPost{}, gorm.ErrRecordNotFound
	}

	return content, nil

	// err := r.db.First(nil, id).Error

	// return err == gorm.ErrRecordNotFound, err
}

// TODO: Add Update
// TODO: Add Delete
// TODO: Add RetrieveAll
// TODO: Add Find? Retrieve with filters could be cool (filter by id, author, etc.)
