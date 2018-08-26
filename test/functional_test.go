package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HTTPMethod string

const (
	Get    HTTPMethod = "GET"
	Post   HTTPMethod = "POST"
	Put    HTTPMethod = "PUT"
	Delete HTTPMethod = "DELETE"
)

const APIURL string = "http://0.0.0.0:4242"

var adminToken string
var nonAdminToken string

func TestRegister(t *testing.T) {
	tests := []struct {
		description string

		body    []byte
		isAdmin bool

		expectedCode  int
		expectedError error
	}{
		{
			description: "register non admin - valid",

			body: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),

			expectedCode: 201,
		},
		{
			description: "register admin - valid",

			body: []byte(`
				{
					"email": "bob-admin@vance-refrigeration.com",
					"password": "refrigerator2000",
					"is_admin": true
				}
			`),

			expectedCode: 201,
		},
		{
			description: "register user that already exists - should fail",

			body: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),

			expectedCode: 500,
		},
		{
			description: "register admin while one has already been setup - should fail",

			body: []byte(`
				{
					"email": "michael.scarn@midnight.org",
					"password": "goldenface",
					"is_admin": true
				}
			`),

			expectedCode: 403,
		},
	}

	for _, test := range tests {
		client := &http.Client{}

		t.Run(test.description, func(t *testing.T) {
			req, err := http.NewRequest("POST", APIURL+"/register", bytes.NewReader(test.body))
			if err != nil {
				assert.FailNowf(t, err.Error(), "could not prepare HTTP request for %s%s", APIURL, "/register")
			}

			req.Header.Add("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
			} else {
				assert.Equal(t, test.expectedCode, response.StatusCode, "invalid http code received")
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		description string

		body    []byte
		isAdmin bool

		adminToken    bool
		nonAdminToken bool

		expectedCode  int
		expectedError error
	}{
		{
			description: "register non admin - valid",

			body: []byte(`
				{
					"email": "bob@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			nonAdminToken: true,

			expectedCode: 201,
		},
		{
			description: "login admin - valid",

			body: []byte(`
				{
					"email": "bob-admin@vance-refrigeration.com",
					"password": "refrigerator2000"
				}
			`),
			adminToken: true,

			expectedCode: 201,
		},
	}

	for _, test := range tests {
		client := &http.Client{}

		t.Run(test.description, func(t *testing.T) {
			req, err := http.NewRequest("POST", APIURL+"/login", bytes.NewReader(test.body))
			if err != nil {
				assert.FailNowf(t, err.Error(), "could not prepare HTTP request for %s%s", APIURL, "/login")
			}

			req.Header.Add("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
				}
				if test.adminToken {
					adminToken = strings.Replace(string(body), "\"", "", -1)
				}
				if test.nonAdminToken {
					nonAdminToken = strings.Replace(string(body), "\"", "", -1)
				}
				assert.Equal(t, test.expectedCode, response.StatusCode, "invalid http code received")
			}
		})
	}

}

// TODO: Find a way to remake the authorization functional tests
// For now the only way I can imagine is to parse the valid token that was created in the RegisterTests
// Modify it to trigger the various errors we want to replicate
// And encode it again to be able to test it.
// This will take quite some time however, and doesn't seem like an important task at the moment.

// func TestAuthorization(t *testing.T) {
// 	invalidTokenUntrustedSource := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBkOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.ssULUbgQPn6kt69zutaIvpuUajfTJrqEZ0fs0IlPKyc"
// 	invalidTokenWrongIssuedClaimFormat := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9pbnZhbGlkLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjo5NTIxNzYyMTQ2fQ.-EJPGyU7tZ8v_iPNc3p86_F3HiDRheBgFGh48UvMF5U"
// 	invalidTokenUnknownUserID := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDF8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.kDTvj5rjsP5JLle-lx9P-H9eUxf05I2F1NVlIOyW3VA"
// 	invalidTokenWrongSubClaimFormat := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjo0MiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.FOUdEOCLMmRq_UEjnlaDUdiH0DkpgnTCm3bgA_QqVVY"
// 	invalidTokenMissingSubClaim := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.Oon6v2NOKzoZ3PQsCXkYJ6JaFG1mQCB-R3-HeR30EOg"
// 	invalidTokenExpired := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MTUyMTcyNjE0NiwiZXhwIjoxNDIxNzYyMTQ2fQ.k3jFG6Mw4TfybWtIHv7OE1T0EDthhvgIeo6nksv7SNc"
// 	invalidTokenIssuedInTheFuture := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MjUyMTcyNjE0NiwiZXhwIjoyODIxNzYyMTQ2fQ.jZcnZ_YjkitXoSGmeVNNmv9S6a6z6XDMwjwiqcAWQmA"
// 	invalidTokenBadSignature := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImFsZXgrc2FtcGxlQGJsdWVjYW52YXMuaW8iLCJuYW1lIjoiYWxleCtzYW1wbGVAYmx1ZWNhbnZhcy5pbyIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9iMmZjNGViYzAyNzQyNjAxZmIyZDAyMTAyZGIxZmJhYT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmFsLnBuZyIsIm5pY2tuYW1lIjoiYWxleCtzYW1wbGUiLCJhcHBfbWV0YWRhdGEiOnsiYXV0aG9yaXphdGlvbiI6eyJncm91cHMiOltdfX0sImF1dGhvcml6YXRpb24iOnsiZ3JvdXBzIjpbXX0sImdyb3VwcyI6W10sImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjbGllbnRJRCI6ImtieXVGRGlkTExtMjgwTEl3VkZpYXpPcWpPM3R5OEtIIiwidXBkYXRlZF9hdCI6IjIwMTgtMDMtMjJUMTM6NDI6MjAuMDQxWiIsInVzZXJfaWQiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJpZGVudGl0aWVzIjpbeyJ1c2VyX2lkIjoiNTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwicHJvdmlkZXIiOiJhdXRoMCIsImNvbm5lY3Rpb24iOiJVc2VybmFtZS1QYXNzd29yZC1BdXRoZW50aWNhdGlvbiIsImlzU29jaWFsIjpmYWxzZX1dLCJjcmVhdGVkX2F0IjoiMjAxNy0wNy0xOVQwOTozNDo1OC4yMjlaIiwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2QiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImlhdCI6MjUyMTcyNjE0NiwiZXhwIjoyODIxNzYyMTQ2fQ.JCJuJ8FM_RN785q4p-eqdM-VFYFNUmsQrswa7Ik-vBs"

// 	tests := []struct {
// 		description string

// 		route  string
// 		method HTTPMethod
// 		token  string
// 		body   []byte

// 		expectedCode    int
// 		expectedTitle   string
// 		expectedContent string
// 		expectedAuthor  string
// 		expectedID      string
// 		expectedBody    []byte
// 		expectedError   error
// 	}{
// 		{
// 			description: "valid token read",

// 			route:  "/posts/1",
// 			method: Get,
// 			token:  adminToken,

// 			expectedCode: 404,
// 			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
// 		},
// 		{
// 			// no error expected since read isn't protected
// 			description: "invalid token read",

// 			route:  "/posts/1",
// 			method: Get,
// 			token:  invalidTokenUntrustedSource,

// 			expectedCode: 404,
// 			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
// 		},
// 		{
// 			description: "valid token write",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  adminToken,

// 			expectedCode: 404,
// 			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
// 		},
// 		{
// 			description: "valid token not admin write",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  nonAdminToken,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: record not found"}`),
// 		},
// 		{
// 			description: "untrusted source write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenUntrustedSource,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid 'iss' claim"}`),
// 		},
// 		{
// 			description: "bad signature write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenBadSignature,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid token: signature is invalid"}`),
// 		},
// 		{
// 			description: "expired token write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenExpired,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid claims: Token is expired"}`),
// 		},
// 		{
// 			description: "token issued in the future write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenIssuedInTheFuture,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid claims: Token used before issued"}`),
// 		},
// 		{
// 			description: "missing sub claim write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenMissingSubClaim,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid 'sub' claim: missing claim"}`),
// 		},
// 		{
// 			description: "unknown user id write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenUnknownUserID,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: record not found"}`),
// 		},
// 		{
// 			description: "wrong issued claim write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenWrongIssuedClaimFormat,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid 'iss' claim"}`),
// 		},
// 		{
// 			description: "wrong sub claim write attempt",

// 			route:  "/posts/1",
// 			method: Delete,
// 			token:  invalidTokenWrongSubClaimFormat,

// 			expectedCode: 401,
// 			expectedBody: []byte(`{"message":"could not validate token: invalid 'sub' claim: invalid format"}`),
// 		},
// 	}

// 	for _, test := range tests {
// 		client := &http.Client{}

// 		t.Run(test.description, func(t *testing.T) {
// 			req, err := http.NewRequest(string(test.method), APIURL+test.route, bytes.NewReader(test.body))
// 			if err != nil {
// 				assert.FailNowf(t, err.Error(), "could not prepare HTTP request for %s%s", APIURL, test.route)
// 			}

// 			req.Header.Add("Authorization", "Bearer "+test.token)
// 			req.Header.Add("Content-Type", "application/json")
// 			response, err := client.Do(req)
// 			if err != nil {
// 				assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
// 			} else {
// 				body, err := ioutil.ReadAll(response.Body)
// 				if err != nil {
// 					assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
// 				} else if test.expectedBody != nil {
// 					assert.Equal(t, test.expectedBody, body, "invalid body received")
// 					assert.Equal(t, test.expectedCode, response.StatusCode, "invalid http code received")
// 				}
// 			}
// 		})
// 	}
// }

// TODO: Parse, update and rebuild JWT to be able to compare expected author with the actual author
// Currently since it's randomly generated based on the unix nano time, it's impossible to know what
// to expect.
func TestAPI(t *testing.T) {
	tests := []struct {
		description string

		route  string
		method HTTPMethod
		token  string
		body   []byte

		expectedCode    int
		expectedTitle   string
		expectedContent string
		expectedAuthor  string
		expectedID      string
		expectedBody    []byte
		expectedError   error
	}{
		{
			description: "try to get post that does not exist",

			route:  "/posts/1",
			method: Get,
			token:  adminToken,

			expectedCode: 404,
			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
		},
		{
			description: "try to delete post that does not exist",

			route:  "/posts/1",
			method: Delete,
			token:  adminToken,

			expectedCode: 404,
			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
		},
		{
			description: "try to update post that does not exist",

			route:  "/posts/1",
			method: Put,
			token:  adminToken,
			body: []byte(`
				{
					"title": "ExampleTitle - Postman is great",
					"content": "ExampleContent - It makes it easy to work collaboratively on an API"
				}
			`),

			expectedCode: 404,
			expectedBody: []byte(`{"message":"blog post id 1: resource not found"}`),
		},
		{
			description: "create a post",

			route:  "/posts",
			method: Post,
			token:  adminToken,
			body: []byte(`
				{
					"title": "ExampleTitle - Postman is great",
					"content": "ExampleContent - It makes it easy to work collaboratively on an API"
				}
			`),

			expectedCode:    201,
			expectedID:      `"id":1`,
			expectedTitle:   `"title":"ExampleTitle - Postman is great"`,
			expectedContent: `"content":"ExampleContent - It makes it easy to work collaboratively on an API"`,
			expectedAuthor:  `"author":"bloggo|`,
		},
		{
			description: "get the post that we just created",

			route:  "/posts/1",
			method: Get,
			token:  adminToken,

			expectedCode:    200,
			expectedID:      `"id":1`,
			expectedTitle:   `"title":"ExampleTitle - Postman is great"`,
			expectedContent: `"content":"ExampleContent - It makes it easy to work collaboratively on an API"`,
			expectedAuthor:  `"author":"bloggo|`,
		},
		{
			description: "update the post that we just created",

			route:  "/posts/1",
			method: Put,
			token:  adminToken,
			body: []byte(`
				{
					"title": "Edited title",
					"content": "Edited content"
				}
			`),

			expectedCode: 204,
		},
		{
			description: "get the post that we just updated",

			route:  "/posts/1",
			method: Get,
			token:  adminToken,

			expectedCode:    200,
			expectedID:      `"id":1`,
			expectedTitle:   `"title":"Edited title"`,
			expectedContent: `"content":"Edited content"`,
			expectedAuthor:  `"author":"bloggo|`,
		},
		{
			description: "delete the post",

			route:  "/posts/1",
			method: Delete,
			token:  adminToken,

			expectedCode: 204,
		},
	}

	for _, test := range tests {
		client := &http.Client{}

		t.Run(test.description, func(t *testing.T) {
			req, err := http.NewRequest(string(test.method), APIURL+test.route, bytes.NewReader(test.body))
			if err != nil {
				assert.FailNowf(t, err.Error(), "could not prepare HTTP request for %s%s", APIURL, test.route)
			}

			req.Header.Add("Authorization", "Bearer "+test.token)
			req.Header.Add("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					assert.Equal(t, test.expectedError.Error(), err.Error(), "invalid error received")
				} else if test.expectedBody != nil {
					assert.Equal(t, test.expectedBody, body, "invalid body received")
					assert.Equal(t, test.expectedCode, response.StatusCode, "invalid http code received")
				} else {
					assert.Contains(t, string(body), test.expectedID, "invalid author in body")
					assert.Contains(t, string(body), test.expectedTitle, "invalid title in body")
					assert.Contains(t, string(body), test.expectedContent, "invalid content in body")
					assert.Contains(t, string(body), test.expectedAuthor, "invalid author in body")
					assert.Equal(t, test.expectedCode, response.StatusCode, "invalid http code received")
				}
			}
		})
	}
}
