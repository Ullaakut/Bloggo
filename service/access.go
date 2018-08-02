package service

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// UserRepository reprensents a user repository to get users from their ids
type UserRepository interface {
	Retrieve(id string) (bool, error)
}

// Access is a service that verifies access tokens
type Access struct {
	users         UserRepository
	trustedSource string

	log *zerolog.Logger
}

// NewAccess creates and configures an Access service
func NewAccess(log *zerolog.Logger, userRepository UserRepository, expectedSource string) *Access {
	return &Access{
		log:           log,
		users:         userRepository,
		trustedSource: expectedSource,
	}
}

// ValidateToken decodes the user info in an ID token and checks its iss, sub and exp claims
func (a *Access) ValidateToken(IDToken string) (string, error) {
	// Since the token's signature is already verified when getting the token information,
	// We can skil the verifications in the parser and use it just to parse the token
	p := &jwt.Parser{
		SkipClaimsValidation: true,
	}
	token, _, err := p.ParseUnverified(IDToken, jwt.MapClaims{})
	if err != nil {
		return "", errors.Wrap(err, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token: can't parse claims")
	}

	// Verifies if token is expired or not yet valid (exp claim)
	err = claims.Valid()
	if err != nil {
		return "", errors.Wrap(err, "invalid claims")
	}

	// Verifies the iss claim
	if !claims.VerifyIssuer(a.trustedSource, true) {
		return "", errors.New("invalid 'iss' claim")
	}

	// Verifies the sub claim
	userID, err := a.verifySubject(claims)
	if err != nil {
		return "", errors.Wrap(err, "invalid 'sub' claim")
	}

	a.log.Debug().Str("user_id", userID).Msg("extracted user id from token")

	isAdmin, err := a.users.Retrieve(userID)
	if err != nil {
		return userID, err
	}
	if !isAdmin {
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
