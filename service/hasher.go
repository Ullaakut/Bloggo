package service

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher describes something that can hash passwords
type Hasher interface {
	Hash(password string) (string, error)
}

// HashComparer describes something that can compare hashed passwords with clear passwords
type HashComparer interface {
	Compare(hash, password string) error
}

// HasherComparer describes something that can hash passwords and compare them with clear passwords
type HasherComparer interface {
	Hasher
	HashComparer
}

// BcryptHasher implements the Hasher and HashComparer interfaces and uses Provos and
// Mazières's bcrypt adaptive hashing algorithm
type BcryptHasher struct {
	runs int
}

// NewBcryptHasher instanciates a BcryptHasher and sets its number of runs
func NewBcryptHasher(runs int) HasherComparer {
	return &BcryptHasher{
		runs: runs,
	}
}

// Hash hashes a password
func (bh *BcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bh.runs)
	return string(hashedPassword), err
}

// Compare compares a hashed password to a clear password
func (bh *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
