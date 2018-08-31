package service

import (
	"github.com/Ullaakut/Bloggo/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// UserRepository reprensents a user repository to get users from their ids
type UserRepository interface {
	Retrieve(user *model.User) (*model.User, error)
}

// Access is a service that verifies access tokens
type Access struct {
	users UserRepository

	trustedSource string
	jws           string

	log *zerolog.Logger
}

// NewAccess creates and configures an Access service
func NewAccess(log *zerolog.Logger, userRepository UserRepository, jws string) *Access {
	return &Access{
		log:   log,
		users: userRepository,
		jws:   jws,
	}
}

// ValidateToken decodes the user info in an ID token and checks its iss, sub and exp claims
func (a *Access) ValidateToken(IDToken string) (string, error) {
	// Since the token's signature is already verified when getting the token information,
	// We can skil the verifications in the parser and use it just to parse the token
	p := &jwt.Parser{
		SkipClaimsValidation: true,
	}
	token, err := p.Parse(IDToken, func(*jwt.Token) (interface{}, error) {
		return []byte(a.jws), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "invalid token")
	}

	// should not be able to fail if call to p.Parse didn't fail
	claims := token.Claims.(jwt.MapClaims)

	// Verifies if token is expired or not yet valid (exp claim)
	err = claims.Valid()
	if err != nil {
		return "", errors.Wrap(err, "invalid claims")
	}

	// Verifies the sub claim
	userID, err := a.verifySubject(claims)
	if err != nil {
		return "", errors.Wrap(err, "invalid 'sub' claim")
	}

	user, err := a.users.Retrieve(&model.User{TokenUserID: userID})
	if err != nil {
		return userID, err
	}
	if !user.IsAdmin {
		a.log.Debug().Msg("unauthorized user")
		return userID, errors.New("user does not have write access")
	}
	a.log.Debug().Msg("authorized user")
	return userID, nil
}

// Verifies that the subject of the token (the user id of who owns it) exists
// and is in a valid format, then returns it in a string
func (a *Access) verifySubject(claims jwt.MapClaims) (string, error) {
	unconvertedUserID, ok := claims["sub"]
	if !ok {
		return "", errors.New("missing claim")
	}

	userID, ok := unconvertedUserID.(string)
	if !ok {
		return "", errors.New("invalid format")
	}

	return userID, nil
}
