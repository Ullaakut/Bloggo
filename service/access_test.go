package service

import (
	"testing"

	"github.com/Ullaakut/Bloggo/repo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewAccess(t *testing.T) {
	iss := "https://samples.auth0.com/"

	userRepositoryMock := &repo.UserRepositoryMock{}

	a := NewAccess(userRepositoryMock, iss)

	assert.Equal(t, iss, a.trustedSource, "unexpected issuer set")
	assert.Equal(t, userRepositoryMock, a.users, "unexpected user repo set")
}

func TestValidateToken(t *testing.T) {
	validSub := "auth0|596f27c2c3709661e9cea37d"
	invalidSub := "auth1|596f27c2c3709661e9cea37d"

	testCases := []struct {
		description string

		token         string
		repositoryErr error
		tokenErr      bool // used to know whether or not code will reach repo call

		expectedError error
	}{
		// TODO: Make user repository return that the user exists
		{
			description: "valid token, no errors",

			token: `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`,
		},
		{
			description: "invalid token, source not trusted",

			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBkOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.CrwwiC0MB_xtJCii1jDo8prLXPe3f8WkFOtth-F6-90",
			tokenErr: true,

			expectedError: errors.New("invalid 'iss' claim"),
		},
		{
			description: "invalid token, 'iss' claim format is wrong",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9pbnZhbGlkLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjo5NTIxNzYyMTQ2fQ.JHPopJXgHOzEJzHSu7xF224pY_1ZzAdwz9-wNkpYDHI`,
			tokenErr: true,

			expectedError: errors.New("invalid 'iss' claim"),
		},
		{
			description: "invalid token, user id in sub claim doesnt exist",

			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDF8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.x0l5UqIH7R6cGpAuRmBv8c_UIC_FXRHOlR3-qz6czBs",

			expectedError: errors.New("user not found"),
		},
		{
			description: "invalid token, sub claim is a number instead of a string",

			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjo0MiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.LTc5XsVNIzJ8NJ2B_VxipjshKjFYV-FtQuma4YZB_Tw",
			tokenErr: true,

			expectedError: errors.New("invalid 'sub' claim"),
		},
		{
			description: "invalid token format",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.sdasadas.77LWIWWj9bvH92uLs05BGB61-ikr3vawwgH5RyuKevg`,
			tokenErr: true,

			expectedError: errors.New("invalid character '±' looking for beginning of value"),
		},
		{
			description: "invalid token, expired",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjoxNDIxNzYyMTQ2fQ.t0Kr3EAUbKoC201NLFoQGtGCFNxSytQqYaFE1uNjKSs`,
			tokenErr: true,

			expectedError: errors.New("Token is expired"),
		},
		{
			description: "invalid token, issued in the future",

			token:    `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MjUyMTcyNjE0NiwiZXhwIjoyODIxNzYyMTQ2fQ.JCJuJ8FM_RN785q4p-eqdM-VFYFNUmsQrswa7Ik-vBs`,
			tokenErr: true,

			expectedError: errors.New("Token used before issued"),
		},
		// TODO: Add test where user repository returns that the user doesnt exist
		{
			description: "invalid token, repository error",

			token:         `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjo5NTIxNzcyMTQ2fQ.xLS0znCvhO6-s2Tx_AfdQ3IvjVXrHXIsMY-ifjwHfCM`,
			repositoryErr: errors.New("user not found"),
			tokenErr:      false,

			expectedError: errors.New("user not found"),
		},
	}
	for idx, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			userRepositoryMock := &repo.UserRepositoryMock{}

			a := &Access{
				users:         userRepositoryMock,
				trustedSource: "https://samples.auth0.com/",
			}

			if !testCase.tokenErr {
				userRepositoryMock.On("Retrieve", validSub).Return(testCase.repositoryErr)
				userRepositoryMock.On("Retrieve", invalidSub).Return(errors.New("user not found"))
			}

			err := a.ValidateToken(testCase.token)

			if testCase.expectedError != nil {
				assert.NotEqual(t, nil, err, "unexpected success in test case %d", idx)
				assert.Equal(t, testCase.expectedError.Error(), err.Error(), "wrong error returned in test case %d", idx)
			} else {
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

	a := NewAccess(
		userRepositoryMock,
		"https://samples.auth0.com/",
	)

	for n := 0; n < b.N; n++ {
		a.ValidateToken(token)
	}
}
