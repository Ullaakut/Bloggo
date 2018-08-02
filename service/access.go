package service

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserRepository reprensents a user repository to get users from their ids
type UserRepository interface {
	Retrieve(id string) error
}

// Access is a service that verifies access tokens
type Access struct {
	users         UserRepository
	trustedSource string
}

// NewAccess creates and configures an Access service
func NewAccess(userRepository UserRepository, expectedSource string) *Access {
	return &Access{
		users:         userRepository,
		trustedSource: expectedSource,
	}
}

// ValidateToken decodes the user info in an ID token and checks its iss, sub and exp claims
func (a *Access) ValidateToken(IDToken string) error {
	// Since the token's signature is already verified when getting the token information,
	// We can skil the verifications in the parser and use it just to parse the token
	p := &jwt.Parser{
		SkipClaimsValidation: true,
	}
	token, _, err := p.ParseUnverified(IDToken, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token: can't parse claims")
	}

	// Verifies if token is expired or not yet valid (exp claim)
	err = claims.Valid()
	if err != nil {
		return err
	}

	// Verifies the iss claim
	if !claims.VerifyIssuer(a.trustedSource, true) {
		return errors.New("invalid 'iss' claim")
	}

	// Verifies the sub claim
	userID, err := a.verifySubject(claims)
	if err != nil {
		return err
	}

	return a.users.Retrieve(userID)
}

// Verifies that the subject of the token (the user id of who owns it) exists
// and is in a valid format, then returns it in a string
func (a *Access) verifySubject(claims jwt.MapClaims) (string, error) {
	unconvertedUserID, ok := claims["sub"]
	if !ok {
		return "", errors.New("missing 'sub' claim")
	}

	userID, ok := unconvertedUserID.(string)
	if !ok {
		return "", errors.New("invalid 'sub' claim")
	}

	return userID, nil
}
