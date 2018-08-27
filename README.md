# Bloggo

Bloggo is a blog CRUD that uses JWT tokens for authorization.

<p align="center">
    <img src="images/bloggologo.png" width="350"/>
</p>
<p align="center">
    <a href="#license">
        <img src="https://img.shields.io/badge/license-Apache-blue.svg?style=flat" />
    </a>
    <a href="https://goreportcard.com/report/github.com/Ullaakut/Bloggo">
        <img src="https://goreportcard.com/badge/github.com/Ullaakut/Bloggo" />
    </a>
    <a href="https://github.com/Ullaakut/Bloggo/releases/latest">
        <img src="https://img.shields.io/github/release/Ullaakut/Bloggo.svg?style=flat" />
    </a>
    <a href="https://travis-ci.org/Ullaakut/Bloggo">
        <img src="https://travis-ci.org/Ullaakut/Bloggo.svg?branch=master" />
    </a>
    <a href='https://coveralls.io/github/Ullaakut/Bloggo?branch=add-ci'>
        <img src='https://coveralls.io/repos/github/Ullaakut/Bloggo/badge.svg?branch=add-ci' alt='Coverage Status' />
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

<p align="center">
    <img width="500"  src="images/blueprints.png">
</p>

## Postman collection

* Import the [collection](/postman/Bloggo.postman_collection.json) in Postman.

<p align="center">
    <img width="500" src="images/postmanCollection.png">
    <img width="800" style="margin-bottom:55px" src="images/postmanExample.png">
</p>

## Testing

To test the API, I recommend [importing the postman collection and using Postman](#postman-collection).

You will need to first register a new user using one of the provided requests, and to update the collection's authentication token with the token the registration API returned you.

In order to do that, hover the `Bloggo` collection name in Postman and click on the three dots to see the collection options. Then, click on `Edit` and navigate to the `Authorization` tab. Finally, paste your token in the `Token` field, without the double quotes.

<p align="center">
    <img width="800" src="images/postmanUpdateToken.png">
</p>

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
