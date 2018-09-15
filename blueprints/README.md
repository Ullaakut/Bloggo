# API documentation

The blueprints represent the API's documentation.

They document for each route:

* The supported methods
* The expected inputs
* The API's responses

## How to build the blueprints

The simplest way is to use the docker image, either by running `docker build -t . blueprints` or by using the `docker-compose.yml` file at the root of the repository and running `docker-compose up blueprints`

You can also build them yourself by installing `nodejs`, and running `npm install -g aglio`, and then running `aglio -h 0.0.0.0 -i main.apib --theme-variables slate -s`

## How to access the blueprints

Once they are deployed, you can find the blueprints in your browser by accessing `http://0.0.0.0:3000`

## Screenshots

<p align="center">
    <img width="500"  src="../images/blueprints.png">
</p>
