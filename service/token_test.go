package service

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Ullaakut/Bloggo/logger"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/Ullaakut/Bloggo/repo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	jws := "MySuperSecretSecret"

	userRepositoryMock := &repo.UserRepositoryMock{}

	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	a := NewToken(log, userRepositoryMock, jws)

	assert.Equal(t, jws, a.jws, "unexpected jws set")
	assert.Equal(t, log, a.log, "unexpected logger set")
	assert.Equal(t, userRepositoryMock, a.user, "unexpected user repo set")
}

func TestGenerateID(t *testing.T) {
	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	a := &Token{
		log: log,
	}

	id := a.GenerateID()
	assert.Contains(t, id, "bloggo|", "unexpected ID generated")
}

func TestLogin(t *testing.T) {
	tests := []struct {
		description string

		userInfo   *model.User
		actualUser *model.User
		repoError  error

		expectedFirstSegment string
		expectedThirdSegment string
		expectedError        error
	}{
		{
			description: "valid token, no errors",

			userInfo: &model.User{
				Email:    "bob@vance-refrigeration.com",
				Password: "refrigerator2000",
			},
			actualUser: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
			},

			expectedFirstSegment: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			// Can't verify the second and third segments without faking the time.Now() call
			expectedError: nil,
		},
		{
			description: "wrong password",

			userInfo: &model.User{
				Email:    "bob@vance-refrigeration.com",
				Password: "wrong password",
			},
			actualUser: &model.User{
				Email:       "bob@vance-refrigeration.com",
				Password:    "refrigerator2000",
				TokenUserID: "test",
			},

			expectedError: errors.New("invalid password"),
		},
		{
			description: "user does not exist",

			userInfo: &model.User{
				Email:    "bob@vance-refrigeration.com",
				Password: "refrigerator2000",
			},
			actualUser: nil,
			repoError:  errors.New("dummy error"),

			expectedError: errors.New("user not found: dummy error"),
		},
	}

	for idx, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			userRepositoryMock := &repo.UserRepositoryMock{}
			userRepositoryMock.
				On("Retrieve", &model.User{Email: test.userInfo.Email}).
				Return(test.actualUser, test.repoError).
				Once()

			a := &Token{
				log:  log,
				jws:  "x5fVmkmyMLAQJiJ8rvsGEAgetl9GS7j8",
				user: userRepositoryMock,
			}

			token, err := a.Login(test.userInfo)

			if test.expectedError != nil {
				assert.NotEqual(t, nil, err, "unexpected success in test case %d", idx)
				assert.Equal(t, test.expectedError.Error(), err.Error(), "wrong error returned in test case %d", idx)
			} else {
				segments := strings.Split(token, ".")
				assert.Equal(t, test.expectedFirstSegment, segments[0], "unexpected token in test case %d", idx)
				assert.Equal(t, nil, err, "unexpected error in test case %d", idx)
			}
		})
	}
}
