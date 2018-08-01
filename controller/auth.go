package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// AccessService represents a service to verifying access tokens
type AccessService interface {
	ValidateToken(IDToken string) error
}

// Auth is a controller that is in charge of authenticating and authorizing requests
type Auth struct {
	access AccessService
}

// NewAuth creates an Auth controller with the given endpoint
func NewAuth(access AccessService) *Auth {
	return &Auth{
		access: access,
	}
}

// Authorize returns the user info for the given access token in the Authorization header
// It is a middleware, which is why it returns a function that, if successful, calls the
// function bound to the route.
func (a *Auth) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Parse Auth header
		token, err := parseAuth(ctx.Request().Header.Get("Authorization"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprint("could not parse auth header: ", err))
		}

		// Verify token claims and expiration date
		err = a.access.ValidateToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprint("could not validate token: ", err))
		}

		return next(ctx)
	}
}

func parseAuth(auth string) (string, error) {
	// check if authorization header exists
	if len(auth) == 0 {
		return "", errors.New("missing authorization header")
	}

	// check if authorization header is valid
	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "", errors.Errorf("invalid authorization header format (%v)", auth)
	}

	// check if we have a bearer token
	if parts[0] != "Bearer" {
		return "", errors.Errorf("invalid authorization header type (%v)", parts[0])
	}

	return parts[1], nil
}
