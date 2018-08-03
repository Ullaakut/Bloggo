package controller

import (
	"net/http"
	"strconv"
	"time"

	errortypes "github.com/Ullaakut/Bloggo/errorTypes"
	"github.com/Ullaakut/Bloggo/model"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	v "gopkg.in/go-playground/validator.v9"
)

// BlogRepository represents a repository that allows to create, read, update and delete blog posts
type BlogRepository interface {
	Retrieve(id uint) (*model.BlogPost, error)
	RetrieveAll() ([]*model.BlogPost, error)
	Store(post *model.BlogPost) (*model.BlogPost, error)
}

// Blog is a controller that is in charge of handling the CRUD of blog posts
type Blog struct {
	posts BlogRepository

	log *zerolog.Logger
}

// NewBlog creates a Blog controller with the given blog post repository
func NewBlog(log *zerolog.Logger, blogPostRepository BlogRepository) *Blog {
	return &Blog{
		posts: blogPostRepository,

		log: log,
	}
}

// Create creates a new blog post
func (b *Blog) Create(ctx echo.Context) error {
	var post model.BlogPost

	err := ctx.Bind(&post)
	if err != nil {
		err = errors.Wrap(err, "could not parse blog post from request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, ok := ctx.Get("userID").(string)
	if !ok {
		err := errors.New("userID not set in request context")
		b.log.Error().Err(err).Msg("potential security breach")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	validate := v.New()
	err = validate.Struct(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Set the author to the user ID so that the API can't be used manually
	// to claim that a post was created by another user
	post.Author = userID

	// Set createdAt value to now
	post.CreatedAt = time.Now()

	createdPost, err := b.posts.Store(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, createdPost)
}

// Read retrieves a blog post from its id
func (b *Blog) Read(ctx echo.Context) error {
	// parse the ID from the URL parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "could not parse blog post ID")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// retrieve the blog post from the blog post repository
	blogPost, err := b.posts.Retrieve(uint(id))
	if errors.Cause(err) == errortypes.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, errors.Wrapf(err, "blog post id %d", id))
	}
	if err != nil {
		err = errors.Wrap(err, "could not read blog post")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, blogPost)
}

// ReadAll retrieves all blog posts
func (b *Blog) ReadAll(ctx echo.Context) error {
	blogPosts, err := b.posts.RetrieveAll()
	if err != nil {
		err = errors.Wrap(err, "could not read blog posts")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, blogPosts)
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
