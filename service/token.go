package service

import (
	"fmt"
	"time"

	"github.com/Ullaakut/Bloggo/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Token is a service that generates JWT tokens
type Token struct {
	jws string

	user UserRepository
	hash HashComparer

	log *zerolog.Logger
}

// NewToken creates and configures an Token service
func NewToken(log *zerolog.Logger, user UserRepository, hash HashComparer, jws string) *Token {
	return &Token{
		log:  log,
		user: user,
		hash: hash,
		jws:  jws,
	}
}

// GenerateID generate a unique ID
func (t *Token) GenerateID() string {
	return "bloggo|" + jwt.EncodeSegment([]byte(fmt.Sprint(time.Now().UnixNano())))
}

// Login generates a signed JWT from the user information if it's valid
func (t *Token) Login(userInfo *model.User) (string, error) {

	actualUser, err := t.user.Retrieve(&model.User{Email: userInfo.Email})
	if err != nil {
		return "", errors.Wrap(err, "user not found")
	}

	err = t.hash.Compare(actualUser.Password, userInfo.Password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   actualUser.TokenUserID,
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.jws))
}
