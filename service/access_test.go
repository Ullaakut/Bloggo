package service

import (
	"bytes"
	"testing"

	"github.com/Ullaakut/Bloggo/logger"
	"github.com/Ullaakut/Bloggo/model"
	"github.com/Ullaakut/Bloggo/repo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewAccess(t *testing.T) {
	jws := "https://samples.auth0.com/"

	userRepositoryMock := &repo.UserRepositoryMock{}

	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	a := NewAccess(log, userRepositoryMock, jws)

	assert.Equal(t, jws, a.jws, "unexpected jws set")
	assert.Equal(t, userRepositoryMock, a.users, "unexpected user repo set")
}

func TestValidateToken(t *testing.T) {
	validSub := "auth0|596f27c2c3709661e9cea37d"
	invalidSub := "auth1|596f27c2c3709661e9cea37d"
	notAdmin := "auth2|596f27c2c3709661e9cea37d"

	tests := []struct {
		description string

		token         string
		repositoryErr error
		tokenErr      bool // used to know whether or not code will reach repo call

		expectedUserID string
		expectedError  error
	}{
		{
			description: "valid token, no errors",

			token: `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.E9z_zYD66k1CnfDVLe16oVuK_c-JYB0g2B0HbJ-PQcA`,

			expectedUserID: "auth0|596f27c2c3709661e9cea37d",
		},
		{
			description: "invalid token, user id in sub claim doesnt exist",

			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDF8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.kDTvj5rjsP5JLle-lx9P-H9eUxf05I2F1NVlIOyW3VA",

			expectedError: errors.New("user not found"),
		},
		{
			description: "invalid token, sub claim is a number instead of a string",

			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjo0MiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.FOUdEOCLMmRq_UEjnlaDUdiH0DkpgnTCm3bgA_QqVVY",
			tokenErr: true,

			expectedError: errors.New("invalid 'sub' claim: invalid format"),
		},
		{
			description: "invalid token, missing 'sub' claim",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.Oon6v2NOKzoZ3PQsCXkYJ6JaFG1mQCB-R3-HeR30EOg`,
			tokenErr: true,

			expectedError: errors.New("invalid 'sub' claim: missing claim"),
		},
		{
			description: "invalid token, expired",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjoxNDIxNzYyMTQ2fQ.k3jFG6Mw4TfybWtIHv7OE1T0EDthhvgIeo6nksv7SNc`,
			tokenErr: true,

			expectedError: errors.New("invalid claims: Token is expired"),
		},
		{
			description: "invalid token, issued in the future",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MjUyMTcyNjE0NiwiZXhwIjoyODIxNzYyMTQ2fQ.jZcnZ_YjkitXoSGmeVNNmv9S6a6z6XDMwjwiqcAWQmA`,
			tokenErr: true,

			expectedError: errors.New("invalid claims: Token used before issued"),
		},
		{
			description: "invalid token: signature is invalid",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MjUyMTcyNjE0NiwiZXhwIjoyODIxNzYyMTQ2fQ.JCJuJ8FM_RN785q4p-eqdM-VFYFNUmsQrswa7Ik-vBs`,
			tokenErr: true,

			expectedError: errors.New("invalid token: signature is invalid"),
		},
		{
			description: "valid token, but not an admin",

			token: `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDJ8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.NJ2zehMX7bgkSG1-Rp0_93rMI0kzVO-2OlzV4FyfoNg`,

			expectedUserID: "auth2|596f27c2c3709661e9cea37d",
			expectedError:  errors.New("user does not have write access"),
		},
		{
			description: "invalid token, repository error",

			token:         `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.E9z_zYD66k1CnfDVLe16oVuK_c-JYB0g2B0HbJ-PQcA`,
			repositoryErr: errors.New("user not found"),
			tokenErr:      false,

			expectedError: errors.New("user not found"),
		},
	}

	for idx, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			userRepositoryMock := &repo.UserRepositoryMock{}

			logsBuff := &bytes.Buffer{}
			log := logger.NewZeroLog(logsBuff)

			a := &Access{
				log:           log,
				users:         userRepositoryMock,
				trustedSource: "https://samples.auth0.com/",
				jws:           "x5fVmkmyMLAQJiJ8rvsGEAgetl9GS7j8",
			}

			if !test.tokenErr {
				userRepositoryMock.On("Retrieve", &model.User{TokenUserID: validSub}).Return(&model.User{IsAdmin: true}, test.repositoryErr)
				userRepositoryMock.On("Retrieve", &model.User{TokenUserID: notAdmin}).Return(&model.User{IsAdmin: false}, nil)
				userRepositoryMock.On("Retrieve", &model.User{TokenUserID: invalidSub}).Return(nil, errors.New("user not found"))
			}

			userID, err := a.ValidateToken(test.token)

			if test.expectedError != nil {
				assert.NotEqual(t, nil, err, "unexpected success in test case %d", idx)
				assert.Equal(t, test.expectedError.Error(), err.Error(), "wrong error returned in test case %d", idx)
			} else {
				assert.Equal(t, test.expectedUserID, userID, "unexpected userID in test case %d", idx)
				assert.Equal(t, nil, err, "unexpected error in test case %d", idx)
			}
		})
	}
}

// BenchmarkValidateToken benchmarks the token validation method
// 3702ns per op on average on a 15" MBP 2017
func BenchmarkValidateToken(b *testing.B) {
	token := `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjoxNTIxNzYyMTQ2fQ.cLRLwTtNAcF4nVo0kOxPL7SfJFkku8IWWx8B812txy8`
	userRepositoryMock := &repo.UserRepositoryMock{}

	logsBuff := &bytes.Buffer{}
	log := logger.NewZeroLog(logsBuff)

	a := &Access{
		log:           log,
		users:         userRepositoryMock,
		trustedSource: "https://samples.auth0.com/",
	}

	for n := 0; n < b.N; n++ {
		a.ValidateToken(token)
	}
}
