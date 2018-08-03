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

To get a valid admin account, use this token:

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`

To get an account that isn't admin, either connect using your own credentials, or use those credentials:

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYytpbnZhbGlkaWRAZXBpdGVjaC5ldSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2UiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImV4cCI6MTYwMDQ5Mjk2NSwiaWF0IjoxNTAwNDU2OTY1fQ.-a86RgUnCqVlplFff6-44E-FejFMuvWK5qEzLoUywhU`

To use your own account on Bloggo, simply get a JWT token from [OpenIDConnect](https://openidconnect.net/), and use adminer to add your user ID to the admins like such:

<img width="300" alt="screenshot 2018-08-03 at 23 18 44" src="https://user-images.githubusercontent.com/6976628/43666375-ee4698b0-9773-11e8-9b23-103fad88cf4b.png">

Then, you should be able to use your token on the API. Don't forget to add `Bearer` before your token in the `Authorization` header though.

## Possible future improvements

See the [issues](https://github.com/Ullaakut/Bloggo/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Potential+future+improvements%22) and [projects](https://github.com/Ullaakut/Blogger/projects/2) pages for a list of possible improvements that could be done in the next 14 days.

### Example of a valid token

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`

### Notes

Echo actually provides a [JWT middleware](https://echo.labstack.com/middleware/jwt), but I feel like it would be cheating considering this is a coding exercise.
