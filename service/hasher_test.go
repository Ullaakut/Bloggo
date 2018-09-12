package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBcryptHasher(t *testing.T) {
	NewBcryptHasher(42)
}

func TestHash(t *testing.T) {
	tests := []struct {
		description string

		password string

		expectedHashBeginning string
		expectedHashLength    int
		expectedErr           error
	}{
		{
			description: "valid test",

			password: "test",

			expectedHashBeginning: "$2a$11",
			expectedHashLength:    len("$2a$11$99n/E61OMEtfFRsyrvHb6uvB80wrOCCV6zXRxujItx7zh.jJeiFpW"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			hasher := &BcryptHasher{
				runs: 11,
			}

			hash, err := hasher.Hash(test.password)
			if err != nil {
				assert.Equal(t, test.expectedErr, err, "unexpected error")
			} else {
				assert.Equal(t, test.expectedHashBeginning, hash[0:6], "unexpected hash beginning")
				assert.Equal(t, test.expectedHashLength, len(hash), "unexpected hash length")
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		description string

		password string
		hash     string

		expectedErr error
	}{
		{
			description: "valid test",

			password: "test",
			hash:     "$2a$11$tP88VYt33B1vvCnzdm9UO.x/gVTxcB8mUtWma0/ba9YbZ6s51.r2y",

			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			hasher := &BcryptHasher{
				runs: 11,
			}

			err := hasher.Compare(test.hash, test.password)
			assert.Equal(t, test.expectedErr, err, "unexpected error")
		})
	}
}
