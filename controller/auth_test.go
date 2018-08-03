package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ullaakut/Bloggo/logger"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccessMock struct {
	mock.Mock
}

func (m *AccessMock) ValidateToken(IDToken string) (string, error) {
	args := m.Called(IDToken)

	return args.String(0), args.Error(1)
}

func TestNewAuth(t *testing.T) {

	accessMock := &AccessMock{}

	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	a := NewAuth(log, accessMock)

	assert.Equal(t, accessMock, a.access, "unexpected access service set")
	assert.Equal(t, log, a.log, "unexpected logger set")
}

func TestAuthorize(t *testing.T) {
	fakeToken := "fakeToken"

	tests := []struct {
		description string

		authHeader        string
		validAuthHeader   bool
		missingAuthHeader bool

		validClaimsErr error

		expectedHTTPCode int
		expectedHTTPBody []byte
	}{
		{
			description: "valid token & auth header",

			authHeader:      "Bearer fakeToken",
			validAuthHeader: true,

			validClaimsErr: nil,

			expectedHTTPCode: http.StatusOK,
			expectedHTTPBody: []byte("{}"),
		},
		{
			description: "invalid auth header format: no token",

			authHeader:      "Bearer",
			validAuthHeader: false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: invalid authorization header format (Bearer)"),
		},
		{
			description: "invalid auth header format: no bearer",

			authHeader:      "Nothing fakeToken",
			validAuthHeader: false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: invalid authorization header type (Nothing)"),
		},
		{
			description: "missing auth header",

			missingAuthHeader: true,
			validAuthHeader:   false,

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not parse auth header: missing authorization header"),
		},
		{
			description: "access service fails",

			authHeader:      "Bearer fakeToken",
			validAuthHeader: true,

			validClaimsErr: errors.New("dummy error"),

			expectedHTTPCode: http.StatusUnauthorized,
			expectedHTTPBody: []byte("could not validate token"),
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// initialize the echo context to use for the test
			e := echo.New()
			r, err := http.NewRequest(echo.GET, "/", nil)
			if err != nil {
				t.Fatal("could not create request")
			}
			if !test.missingAuthHeader {
				r.Header.Set("Authorization", test.authHeader)
			}

			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			// Setup access service mock
			accessMock := &AccessMock{}
			if test.validAuthHeader {
				accessMock.On("ValidateToken", fakeToken).Return("fakeUserID", test.validClaimsErr).Once()
			}

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			a := Auth{
				log:    log,
				access: accessMock,
			}

			// Since Authorize is a middleware, it needs to be given an HTTP handler to forward
			// the call to, once the authorization is validated. Here we pass it a function that
			// always returns no error :)
			err = a.Authorize(func(ctx echo.Context) error {
				return ctx.JSON(http.StatusOK, struct{}{})
			})(ctx)

			if err == nil {
				assert.Equal(t, test.expectedHTTPCode, w.Code, "wrong response status")
				assert.Equal(t, test.expectedHTTPBody, w.Body.Bytes(), "wrong response body")
			} else {
				assert.Contains(t, err.Error(), fmt.Sprint(test.expectedHTTPCode), "wrong error response status")
				if test.expectedHTTPBody != nil {
					assert.Contains(t, err.Error(), string(test.expectedHTTPBody), "unexpected error response")
				} else {
					assert.Contains(t, w.Body, nil, "unexpected error response")
				}
			}

			accessMock.AssertExpectations(t)
		})
	}
}
