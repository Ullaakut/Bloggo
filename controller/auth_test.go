package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccessMock struct {
	mock.Mock
}

func (m *AccessMock) ValidateToken(IDToken string) error {
	args := m.Called(IDToken)
	return args.Error(0)
}

func TestNewAuth(t *testing.T) {

	accessMock := &AccessMock{}
	a := NewAuth(accessMock)

	assert.Equal(t, accessMock, a.access, "unexpected access service set")
}

func TestGetUserInfo(t *testing.T) {
	fakeToken := "fakeToken"

	testCases := []struct {
		authHeader        string
		validAuthHeader   bool
		missingAuthHeader bool

		validClaimsErr error

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			authHeader:      "Bearer fakeToken",
			validAuthHeader: true,

			validClaimsErr: nil,

			expectedHTTPCode: http.StatusOK,
		},
		{
			authHeader:      "Bearer",
			validAuthHeader: false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: invalid authorization header format (Bearer)"),
		},
		{
			authHeader:      "Nothing fakeToken",
			validAuthHeader: false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: invalid authorization header type (Nothing)"),
		},
		{
			missingAuthHeader: true,
			validAuthHeader:   false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: missing authorization header"),
		},
		{
			authHeader:      "Bearer fakeToken",
			validAuthHeader: true,

			validClaimsErr: errors.New("dummy error"),

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not validate token"),
		},
	}
	for _, testCase := range testCases {
		// initialize the echo context to use for the test
		e := echo.New()
		r, err := http.NewRequest(echo.GET, "/", nil)
		if err != nil {
			t.Fatal("could not create request")
		}
		if !testCase.missingAuthHeader {
			r.Header.Set("Authorization", testCase.authHeader)
		}

		w := httptest.NewRecorder()
		ctx := e.NewContext(r, w)

		// Setup access service mock
		accessMock := &AccessMock{}
		if testCase.validAuthHeader {
			accessMock.On("ValidateToken", fakeToken).Return(testCase.validClaimsErr).Once()
		}

		a := NewAuth(accessMock)

		// Since Authorize is a middleware, it needs to be given an HTTP handler to forward
		// the call to, once the authorization is validated. Here we pass it a function that
		// always returns no error :)
		err = a.Authorize(func(ctx echo.Context) error {
			return ctx.JSON(http.StatusOK, "")
		})(ctx)

		if err == nil {
			assert.Equal(t, testCase.expectedHTTPCode, w.Code, "wrong response status")
			assert.Equal(t, testCase.expectedHTTPBody, w.Body.Bytes(), "wrong response body")
		} else {
			assert.Contains(t, err.Error(), fmt.Sprint(testCase.expectedHTTPCode), "wrong error response status")
			if testCase.expectedHTTPBody != nil {
				assert.Contains(t, w.Body.Bytes(), string(testCase.expectedHTTPBody), "unexpected error response")
			} else {
				assert.Contains(t, w.Body, nil, "unexpected error response")
			}
		}

		accessMock.AssertExpectations(t)
	}
}
