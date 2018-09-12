package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ullaakut/Bloggo/logger"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/Ullaakut/Bloggo/repo"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TokenGeneratorMock struct {
	mock.Mock
}

func (m *TokenGeneratorMock) Login(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *TokenGeneratorMock) GenerateID() string {
	args := m.Called()
	return args.String(0)
}

type HasherMock struct {
	mock.Mock
}

func (m *HasherMock) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestNewUser(t *testing.T) {
	userRepositoryMock := &repo.UserRepositoryMock{}
	hasherMock := &HasherMock{}
	tokenMock := &TokenGeneratorMock{}
	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	b := NewUser(log, userRepositoryMock, tokenMock, hasherMock)

	assert.Equal(t, userRepositoryMock, b.users, "unexpected user repository set")
	assert.Equal(t, tokenMock, b.tokens, "unexpected token service set")
	assert.Equal(t, hasherMock, b.hasher, "unexpected hashing service set")
	assert.Equal(t, log, b.log, "unexpected logger set")
}

func TestRegister(t *testing.T) {
	tests := []struct {
		description string

		requestBody    []byte
		user           *model.User
		adminExists    bool
		generatedToken string
		generatedHash  string

		repositoryErr error
		loginErr      error
		hashErr       error

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			description: "register: passing test",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			user: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
			},
			generatedHash:  "fakeHash",
			generatedToken: "x.y.z",

			expectedHTTPCode: 201,
			expectedHTTPBody: []byte(`"x.y.z"`),
		},
		{
			description: "register admin: passing test",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000",
					"is_admin": true
				}
			`),
			user: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
				IsAdmin:     true,
			},
			generatedHash:  "fakeHash",
			generatedToken: "x.y.z",

			expectedHTTPCode: 201,
			expectedHTTPBody: []byte(`"x.y.z"`),
		},
		{
			description: "register admin: admin has already been setup",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000",
					"is_admin": true
				}
			`),
			generatedHash:  "fakeHash",
			generatedToken: "x.y.z",
			adminExists:    true,

			expectedHTTPCode: 403,
			expectedHTTPBody: []byte(`admin account has already been created`),
		},
		{
			description: "invalid email address",

			requestBody: []byte(`
				{
					"email": "not-an-email-address",
					"password": "refrigerator2000",
					"is_admin": true
				}
			`),

			expectedHTTPCode: 422,
			expectedHTTPBody: []byte(`Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag`),
		},
		{
			description: "invalid password (too short)",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "12345",
					"is_admin": true
				}
			`),

			expectedHTTPCode: 422,
			expectedHTTPBody: []byte(`Key: 'User.Password' Error:Field validation for 'Password' failed on the 'min' tag`),
		},
		{
			description: "hashing error",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "123456789012345",
					"is_admin": true
				}
			`),

			generatedToken: "x.y.z",

			hashErr: errors.New("could not hash"),

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`could not hash`),
		},
		{
			description: "not json",

			requestBody: []byte(`potato`),

			expectedHTTPCode: 400,
			expectedHTTPBody: []byte(`Syntax error: offset=1, error=invalid character 'p' looking for beginning of value`),
		},
		{
			description: "repo error",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			user: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
			},
			repositoryErr:  errors.New("dummy error"),
			generatedHash:  "fakeHash",
			generatedToken: "x.y.z",

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`dummy error`),
		},
		{
			description: "login error",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			user: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
			},
			loginErr:       errors.New("dummy error"),
			generatedHash:  "fakeHash",
			generatedToken: "x.y.z",

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`dummy error`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// initialize the echo context to use for the test
			e := echo.New()
			r, err := http.NewRequest(echo.POST, "/register", bytes.NewReader(test.requestBody))
			if err != nil {
				t.Fatal("could not create request")
			}
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			userRepositoryMock := &repo.UserRepositoryMock{}
			if (test.generatedHash != "" || test.repositoryErr != nil) && !test.adminExists {
				userRepositoryMock.
					On("Store", mock.AnythingOfType("*model.User")).
					Return(test.user, test.repositoryErr).
					Once()
			}

			if strings.Contains(string(test.requestBody), "is_admin") && test.generatedHash != "" {
				userRepositoryMock.On("AdminExists").Return(test.adminExists).Once()
			}

			tokenMock := &TokenGeneratorMock{}
			if test.generatedToken != "" {
				tokenMock.On("GenerateID").Return("test").Once()
			}
			if test.repositoryErr == nil && !test.adminExists && test.generatedHash != "" {
				tokenMock.
					On("Login", mock.AnythingOfType("*model.User")).
					Return("x.y.z", test.loginErr).
					Once()
			}

			hasherMock := &HasherMock{}
			if test.generatedToken != "" {
				hasherMock.
					On("Hash", mock.AnythingOfType("string")).
					Return(test.generatedHash, test.hashErr).
					Once()
			}

			userController := &User{
				users:  userRepositoryMock,
				tokens: tokenMock,
				hasher: hasherMock,

				log: log,
			}

			err = userController.Register(ctx)

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

			userRepositoryMock.AssertExpectations(t)
			tokenMock.AssertExpectations(t)
			hasherMock.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		description string

		requestBody []byte
		loginErr    error
		validUser   bool

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			description: "login: passing test",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			validUser: true,

			expectedHTTPCode: 201,
			expectedHTTPBody: []byte(`"x.y.z"`),
		},
		{
			description: "invalid email address",

			requestBody: []byte(`
				{
					"email": "not-an-email-address",
					"password": "refrigerator2000",
					"is_admin": true
				}
			`),

			expectedHTTPCode: 422,
			expectedHTTPBody: []byte(`Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag`),
		},
		{
			description: "invalid password (too short)",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "12345",
					"is_admin": true
				}
			`),

			expectedHTTPCode: 422,
			expectedHTTPBody: []byte(`Key: 'User.Password' Error:Field validation for 'Password' failed on the 'min' tag`),
		},
		{
			description: "not json",

			requestBody: []byte(`potato`),

			expectedHTTPCode: 400,
			expectedHTTPBody: []byte(`Syntax error: offset=1, error=invalid character 'p' looking for beginning of value`),
		},
		{
			description: "login error",

			requestBody: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			validUser: true,
			loginErr:  errors.New("dummy error"),

			expectedHTTPCode: 500,
			expectedHTTPBody: []byte(`dummy error`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// initialize the echo context to use for the test
			e := echo.New()
			r, err := http.NewRequest(echo.POST, "/login", bytes.NewReader(test.requestBody))
			if err != nil {
				t.Fatal("could not create request")
			}
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			tokenMock := &TokenGeneratorMock{}
			if test.validUser {
				tokenMock.
					On("Login", mock.AnythingOfType("*model.User")).
					Return("x.y.z", test.loginErr).
					Once()
			}

			hasherMock := &HasherMock{}
			hasherMock.
				On("Hash", mock.AnythingOfType("string")).
				Return("ok", nil).
				Once()

			userController := &User{
				tokens: tokenMock,

				log: log,
			}

			err = userController.Login(ctx)

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

			tokenMock.AssertExpectations(t)
		})
	}
}
