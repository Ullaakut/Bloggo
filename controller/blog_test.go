package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ullaakut/Bloggo/errorTypes"
	"github.com/Ullaakut/Bloggo/logger"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/Ullaakut/Bloggo/repo"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a custom error to imitate a gorm.ErrRecordNotFound without actually
// importing gorm here
type ResourceNotFoundErr struct {
}

func (r *ResourceNotFoundErr) Cause() error {
	return errortypes.ErrNotFound
}

func (r *ResourceNotFoundErr) Error() string {
	return errortypes.ErrNotFound.Error()
}

func TestNewBlog(t *testing.T) {
	blogPostRepositoryMock := &repo.BlogPostRepositoryMock{}
	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	b := NewBlog(log, blogPostRepositoryMock)

	assert.Equal(t, blogPostRepositoryMock, b.posts, "unexpected blog post repository set")
	assert.Equal(t, log, b.log, "unexpected logger set")
}

func TestCreate(t *testing.T) {
	tests := []struct {
		description string

		requestBody     []byte
		userIDMissing   bool
		repositoryErr   error
		createdBlogPost *model.BlogPost

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			description: "created: passing test",

			requestBody: []byte(`
				{
					"title": "lorem ipsum",
					"content": "dolor sit amet"
				}
			`),
			createdBlogPost: &model.BlogPost{
				ID:        1,
				Title:     "lorem ipsum",
				Content:   "dolor sit amet",
				Author:    "faketoken",
				CreatedAt: time.Time{},
			},

			expectedHTTPCode: 201,
			expectedHTTPBody: []byte(`{"id":1,"author":"faketoken","title":"lorem ipsum","content":"dolor sit amet","created_at":"0001-01-01T00:00:00Z"}`),
		},
		{
			description: "bad request: bind fails",

			requestBody: []byte(`Not json`),

			expectedHTTPCode: 400,
			expectedHTTPBody: []byte(`Syntax error: offset=1, error=invalid character 'N' looking for beginning of value`),
		},
		{
			description: "internal server error: could not find userID in context",

			requestBody:   []byte(`{}`),
			userIDMissing: true,

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`userID not set in request context`),
		},
		{
			description: "unprocessable entity: invalid blog post",

			requestBody: []byte(`{}`),

			expectedHTTPCode: 422,
			expectedHTTPBody: []byte(`Error:Field validation for 'Title' failed on the 'required' tag`),
		},
		{
			description: "internal server error: repository failure",

			requestBody: []byte(`
				{
					"title": "lorem ipsum",
					"content": "dolor sit amet"
				}
			`),
			repositoryErr:   errors.New("database exploded"),
			createdBlogPost: nil,

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`database exploded`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// initialize the echo context to use for the test
			e := echo.New()
			r, err := http.NewRequest(echo.POST, "/posts", bytes.NewReader(test.requestBody))
			if err != nil {
				t.Fatal("could not create request")
			}
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			if !test.userIDMissing {
				ctx.Set("userID", "fakeToken")
			}

			blogPostRepositoryMock := &repo.BlogPostRepositoryMock{}
			if test.repositoryErr != nil || test.createdBlogPost != nil {
				blogPostRepositoryMock.
					On("Store", mock.Anything).
					Return(test.createdBlogPost, test.repositoryErr).
					Once()
			}

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			blogController := &Blog{
				posts: blogPostRepositoryMock,

				log: log,
			}

			err = blogController.Create(ctx)

			if err == nil {
				assert.Equal(t, test.expectedHTTPCode, w.Code, "wrong response status")
				assert.Equal(t, string(test.expectedHTTPBody), w.Body.String(), "wrong response body")
			} else {
				assert.Contains(t, err.Error(), fmt.Sprint(test.expectedHTTPCode), "wrong error response status")
				if test.expectedHTTPBody != nil {
					assert.Contains(t, err.Error(), string(test.expectedHTTPBody), "unexpected error response")
				} else {
					assert.Contains(t, w.Body, nil, "unexpected error response")
				}
			}

			blogPostRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		description string

		blogPostIDMissing bool
		repositoryErr     error
		retrievedBlogPost *model.BlogPost

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			description: "created: passing test",

			retrievedBlogPost: &model.BlogPost{
				ID:        1,
				Title:     "lorem ipsum",
				Content:   "dolor sit amet",
				Author:    "faketoken",
				CreatedAt: time.Time{},
			},

			expectedHTTPCode: 200,
			expectedHTTPBody: []byte(`{"id":1,"author":"faketoken","title":"lorem ipsum","content":"dolor sit amet","created_at":"0001-01-01T00:00:00Z"}`),
		},
		{
			description: "bad request: missing blog post id",

			blogPostIDMissing: true,

			expectedHTTPCode: 400,
			expectedHTTPBody: []byte(`could not parse blog post ID: strconv.ParseUint: parsing "": invalid syntax`),
		},
		{
			description: "not found: blog post entity doesnt exist",

			repositoryErr:     &ResourceNotFoundErr{},
			retrievedBlogPost: nil,

			expectedHTTPCode: 404,
			expectedHTTPBody: []byte(`blog post id 42: resource not found`),
		},
		{
			description: "internal server error: repository failure",

			repositoryErr:     errors.New("database exploded"),
			retrievedBlogPost: nil,

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`could not read blog post: database exploded`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// initialize the echo context to use for the test
			e := echo.New()
			r, err := http.NewRequest(echo.GET, "/posts/", nil)
			if err != nil {
				t.Fatal("could not create request")
			}

			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			if !test.blogPostIDMissing {
				ctx.SetParamNames("id")
				ctx.SetParamValues("42")
			}

			blogPostRepositoryMock := &repo.BlogPostRepositoryMock{}
			if test.repositoryErr != nil || test.retrievedBlogPost != nil {
				blogPostRepositoryMock.
					On("Retrieve", uint(42)).
					Return(test.retrievedBlogPost, test.repositoryErr).
					Once()
			}

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			blogController := &Blog{
				posts: blogPostRepositoryMock,

				log: log,
			}

			err = blogController.Read(ctx)

			if err == nil {
				assert.Equal(t, test.expectedHTTPCode, w.Code, "wrong response status")
				assert.Equal(t, string(test.expectedHTTPBody), w.Body.String(), "wrong response body")
			} else {
				assert.Contains(t, err.Error(), fmt.Sprint(test.expectedHTTPCode), "wrong error response status")
				if test.expectedHTTPBody != nil {
					assert.Contains(t, err.Error(), string(test.expectedHTTPBody), "unexpected error response")
				} else {
					assert.Contains(t, w.Body, nil, "unexpected error response")
				}
			}

			blogPostRepositoryMock.AssertExpectations(t)
		})
	}
}
