# Bloggo

Bloggo is a blog CRUD that uses JWT tokens for authorization. See [OpenID Connect (OAuth 2.0) identity provider](https://samples.auth0.com/).

<p align="center"><img src="images/bloggologo.png" width="350"/></p>

## Dependencies

* `go` [Download](https://golang.org/dl/)
* `docker` [Download](https://www.docker.com/community-edition)
* `docker-compose` [Download](https://docs.docker.com/compose/install/)
* `dep` [Download](https://github.com/golang/dep)

## Deployment

* `docker-compose up`

## API Blueprints

* `docker-compose up blueprints`
* Visit `0.0.0.0:3000` in your favorite browser

<img width="500" alt="screenshot 2018-08-03 at 20 55 27" src="https://user-images.githubusercontent.com/6976628/43660729-23b3a326-9760-11e8-9c7d-8d425eff6c02.png">

## Postman collection

* Import the [collection](/postman/Bloggo.postman_collection.json) in Postman.

<img width="161" alt="screenshot 2018-08-03 at 20 10 31" src="https://user-images.githubusercontent.com/6976628/43658775-1dd15cf6-975a-11e8-82c0-258732e24ff9.png">

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

### Example of a valid token

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`

### Notes

Echo actually provides a [JWT middleware](https://echo.labstack.com/middleware/jwt), but I feel like it would be cheating considering this is a coding exercise.
