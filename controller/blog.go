package controller

import (
	"github.com/Ullaakut/Bloggo/model"
	"github.com/labstack/echo"
)

// BlogRepository represents a repository that allows to create, read, update and delete blog posts
type BlogRepository interface {
	Retrieve(id uint) (model.BlogPost, error)
	Store(post model.BlogPost) error
}

// Blog is a controller that is in charge of handling the CRUD of blog posts
type Blog struct {
	posts BlogRepository
}

// NewBlog creates a Blog controller with the given blog post repository
func NewBlog(blogPostRepository BlogRepository) *Blog {
	return &Blog{
		posts: blogPostRepository,
	}
}

// Create creates a new blog post
func (b *Blog) Create(ctx echo.Context) error {
	// TODO: Implement this func
	return nil
}

// Read retrieves a blog post from its id
func (b *Blog) Read(ctx echo.Context) error {
	// TODO: Implement this func
	return nil
}

// ReadAll retrieves all blog posts
func (b *Blog) ReadAll(ctx echo.Context) error {
	// TODO: Implement this func
	return nil
}

// Update edits a blog post from its id
func (b *Blog) Update(ctx echo.Context) error {
	// TODO: Implement this func
	return nil
}

// Delete removes a blog post from its id
func (b *Blog) Delete(ctx echo.Context) error {
	// TODO: Implement this func
	return nil
}
