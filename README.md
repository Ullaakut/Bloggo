# Bloggo

Bloggo is a blog CRUD that uses JWT tokens for authorization.

<p align="center">
    <img src="images/bloggologo.png" width="350"/>
</p>
<p align="center">
    <a src="#license">
        <img src="https://img.shields.io/badge/license-Apache-blue.svg?style=flat" />
    </a>
    <a src="https://goreportcard.com/report/github.com/Ullaakut/Bloggo">
        <img src="https://goreportcard.com/badge/github.com/Ullaakut/Bloggo" />
    </a>
    <a src="https://github.com/Ullaakut/Bloggo/releases/latest">
        <img src="https://img.shields.io/github/release/Ullaakut/Bloggo.svg?style=flat" />
    </a>
</p>

## Table of content

* [How to run it](#how-to-run-it)
* [API Blueprints](#api-blueprints)
* [Postman collection](#postman-collection)
* [Testing](#testing)
* [Possible future improvements](#possible-future-improvements)
* [Notes & technical choices](#notes-and-technical-choices)
* [License](#license)

## Dependencies

* `docker` [Download](https://www.docker.com/community-edition)
* `docker-compose` [Download](https://docs.docker.com/compose/install/)

## How to run it

* `docker-compose up`

## API Blueprints

* `docker-compose up blueprints`
* Visit `0.0.0.0:3000` in your favorite browser

<img width="500" alt="screenshot 2018-08-03 at 20 55 27" src="https://user-images.githubusercontent.com/6976628/43660729-23b3a326-9760-11e8-9c7d-8d425eff6c02.png">

## Postman collection

* Import the [collection](/postman/Bloggo.postman_collection.json) in Postman.

<img width="161" alt="screenshot 2018-08-03 at 20 10 31" src="https://user-images.githubusercontent.com/6976628/43658775-1dd15cf6-975a-11e8-82c0-258732e24ff9.png">

## Demonstration

[Simple demonstration video of the behavior of the API](https://youtu.be/c2bogbT6JB4)

## Testing

You can use the credentials below on [OpenID Connect](https://openidconnect.net/) to issue access tokens for testing.

To get a valid admin account, use this token:

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYyt0ZXN0YXBpQGVwaXRlY2guZXUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vc2FtcGxlcy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTk2ZjI3YzJjMzcwOTY2MWU5Y2VhMzdkIiwiYXVkIjoia2J5dUZEaWRMTG0yODBMSXdWRmlhek9xak8zdHk4S0giLCJleHAiOjE2MDA0OTI5NjUsImlhdCI6MTUwMDQ1Njk2NX0.4To8mYgu4pM7J2G5jBTnKWelCTU1U1jo0QOENVp3pOk`

To get an account that isn't admin, either connect using your own credentials, or use those credentials:

`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImJyZW5kYW4ubGUtZ2xhdW5lYytpbnZhbGlkaWRAZXBpdGVjaC5ldSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiaXNzIjoiaHR0cHM6Ly9zYW1wbGVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1OTZmMjdjMmMzNzA5NjYxZTljZWEzN2UiLCJhdWQiOiJrYnl1RkRpZExMbTI4MExJd1ZGaWF6T3FqTzN0eThLSCIsImV4cCI6MTYwMDQ5Mjk2NSwiaWF0IjoxNTAwNDU2OTY1fQ.-a86RgUnCqVlplFff6-44E-FejFMuvWK5qEzLoUywhU`

To use your own account on Bloggo, simply get a JWT token from [OpenIDConnect](https://openidconnect.net/), and use adminer to add your user ID to the admins like such:

<img width="300" alt="screenshot 2018-08-03 at 23 18 44" src="https://user-images.githubusercontent.com/6976628/43666375-ee4698b0-9773-11e8-9b23-103fad88cf4b.png">

Then, you should be able to use your token on the API. Don't forget to add `Bearer` before your token in the `Authorization` header though.

To test the API, I recommend [importing the postman collection and using Postman](#postman-collection).

## Possible future improvements

See the [issues](https://github.com/Ullaakut/Bloggo/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Potential+future+improvements%22) and [projects](https://github.com/Ullaakut/Blogger/projects/2) pages for a list of possible improvements that could be done in the next 14 days.

## Notes and technical choices

### Methodology

* The master branch contains 1 commit for each pull request that was merged, and this commit contains the issue number to link it with the GitHub issues
* Issues are attributed to milestones & GitHub projects
* Issues are divded in two groups: the issues that I wanted to do before the deadline and those that would be possible future improvements.
* All feature pull requests have precise information on the goal of the PR, instructions on how to test them and images / GIFs to show the features in action
* I decided to use API blueprints to document the API, as it's a pretty clean way to represent an API, and it's easy to deploy with the other services.
* I added a postman collection to make it simple & fun to test Bloggo

### Tests

* The services and controllers are 100% unit tested, but the repositories are not.
* Most unit tests are table-driven since it fits pretty well the needs of such functions.
* I unfortunately didn't do TDD when developing Bloggo. I would have liked to, but since I really don't have a lot of experience with it, I'm currently much slower than when developing normally.

### Project structure

* The structure of the code has quite a lot of abstraction. I used interfaces to inject dependencies in the controllers, in order to make testing easy as well as making it possible to, for example, have the choice between multiple different DB connectors. Also, this allows to define in the interface only the functions that the code actually needs, rather than having access to methods that are not needed for the purpose we need the dependency for.
* Controllers are HTTP handlers that have access to services (business logic) and repositories (database connectors).
* Its configuration is currently hard coded. It is part of the future possible improvements to change that, but I judged it wasn't necessary since there is a `docker-compose.yml` files that ensures that it can run on any machine with docker, the way it was meant to be used.
* Bloggo uses structured logging, which makes it easy in the future to plug it into ELK or Greylog for example, and also provides pretty and clean logs with log levels out of the box

### Database

* I decided to use MySQL instead of DynamoDB since I wouldn't have had enough time to integrate it, but I decided to create a ticket to add a DynamoDB repository
* The database is populated the first time you run it, with the contents of the `.sql` files in `/data`
* The credentials for the database are `root`/`root` and the database name is `bloggo`

### Side notes

* Echo actually provides a [JWT middleware](https://echo.labstack.com/middleware/jwt), but I feel like it would be cheating considering this is a coding exercise.
* The only way to add and remove admins is to use the adminer web interface that is deployed with `docker-compose up` and can be accessed from `0.0.0.0:8080`
* I decided to use `dep` over `glide` as it's actually recommended by `glide`'s documentation at the moment.
* The docker image is built using a multi-stage build, in order to reduce the size of the generated Bloggo image as much as possible. Bloggo is 27MB.

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
