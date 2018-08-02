# Bloggo

Bloggo is a blog CRUD that uses JWT tokens for authorization. See [OpenID Connect (OAuth 2.0) identity provider](https://samples.auth0.com/).

<p align="center"><img src="images/bloggologo.png" width="350"/></p>

## Dependencies

* `go` [Download](https://golang.org/dl/)
* `docker` [Download](https://www.docker.com/community-edition)
* `docker-compose` [Download](https://docs.docker.com/compose/install/)
* `dep` [Download](https://github.com/golang/dep)

## Deployment

* `docker-compose up` // TODO: Add Dockerfile & docker-compose file

## API

See [API blueprints](TODO).

## Testing

You can use the credentials below on [OpenID Connect](https://openidconnect.net/) to issue access tokens for testing.

To get a valid admin account, use those credentials:

TODO: Create credentials for an admin account

To get an account that isn't admin, either connect using your own credentials, or use those credentials:

TODO: Create credentials for non-admin account

## Performance

TODO: Estimate API performance
TODO: Add benchmarks + add previous benchmark of token parsing

## Possible future improvements

See the [issues](https://github.com/Ullaakut/Bloggo/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Potential+future+improvements%22) and [projects](https://github.com/Ullaakut/Blogger/projects/2) pages for a list of possible improvements that could be done in the next 14 days.

## Help

### Required Claims

```json
{
    "iss": "https://samples.auth0.com/",
    "sub": "" // TODO: Get admin account and update this
  }
```

### Example of a valid token

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`

### Example of an invalid token

TODO: Add invalid token

### Notes

Echo actually provides a [JWT middleware](https://echo.labstack.com/middleware/jwt), but I feel like it would be cheating considering this is a coding exercise.

Time management:

* Wednesday 01/08: Spent ~3-4 hours building the base structure of the API (from my previous project [PAPI](https://github.com/Ullaakut/PAPI)) and putting ideas into code & TODOs.
