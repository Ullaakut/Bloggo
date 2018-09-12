package controller

import (
	"net/http"

	"github.com/Ullaakut/Bloggo/model"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	v "gopkg.in/go-playground/validator.v9"
)

// UserRepository represents a repository that allows to create users and
// verify whether or not the admin user has already been setup
type UserRepository interface {
	Store(user *model.User) (*model.User, error)
	AdminExists() bool
}

// TokenGenerator represents a service to generate tokens from user info
type TokenGenerator interface {
	Login(user *model.User) (string, error)
	GenerateID() string
}

// Hasher represents a service that hashes passwords
type Hasher interface {
	Hash(password string) (string, error)
}

// User is a controller that is in charge of handling the CRUD of users
type User struct {
	users  UserRepository
	tokens TokenGenerator
	hasher Hasher

	log *zerolog.Logger
}

// NewUser creates a User controller with the given user repository
func NewUser(log *zerolog.Logger, userRepository UserRepository, tokens TokenGenerator, hasher Hasher) *User {
	return &User{
		users:  userRepository,
		tokens: tokens,
		hasher: hasher,

		log: log,
	}
}

// Register creates a new user
func (u *User) Register(ctx echo.Context) error {
	var user model.User

	err := ctx.Bind(&user)
	if err != nil {
		err = errors.Wrap(err, "could not parse user data from request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := v.New()
	err = validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	user.TokenUserID = u.tokens.GenerateID()

	plainTextPwd := user.Password

	user.Password, err = u.hasher.Hash(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Ensure we don't have more than 1 admin user
	if user.IsAdmin && u.users.AdminExists() {
		return echo.NewHTTPError(http.StatusForbidden, errors.New("admin account has already been created"))
	}

	createdUser, err := u.users.Store(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	createdUser.Password = plainTextPwd

	token, err := u.tokens.Login(createdUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, token)
}

// Login gives a token to the user upon providing their credentials
func (u *User) Login(ctx echo.Context) error {
	var user model.User

	err := ctx.Bind(&user)
	if err != nil {
		err = errors.Wrap(err, "could not parse user data from request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := v.New()
	err = validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	token, err := u.tokens.Login(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, token)
}
