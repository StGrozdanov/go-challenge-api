<h1 align="center">Outdoorsy API</h1>

## Project Concept
The assignment is to create a simple REST API written in Golang that serves JSON to the clients, containing available vehicles for rent.

Sorting by parameters and pagination have to be implemented.

There sort parameters can be used in combination and some interesting cases are present from that. What happens when min_price parameter is higher than max_price parameter ? What happens in case of invalid input like negative price and so on.

## Features

* #### Healths endpoint - returns the status of the database and server.

* #### Metrics endpoint - prometeius metrics for overall server behaviour and resource usage.

* #### CORS interceptor

* #### Database reconnect mech - the connections will not always be stable.

* #### Logs coverage

* #### Parameters validation with descriptive error messages as API response

* #### Unit tests - coverage of the base business logic in the rentals.go file.

* #### Basic Integration tests - coverage of the endpoint response statuses.

## Dependencies

* Gin
* Go validator
* sqlx
* goccy go JSON (faster marshal/unmarshal)
* Logrus
* Testify

## Nice to haves in the future

* More validations (there are couple more edge cases that can be considered)
* More tests (both unit and integration)
* Load tests (using locust for example)
* Caching - most of the data seems to be with static nature, we can decrease the DB load and overall performance by using something like Redis

## REST API endpoints

Healths:

* #### GET /healths - check the database and app status.

* #### GET /metrics - check the app metrics

* #### GET /rentals/:id - get a single rental

* #### GET /rentals - get all rentals. Supports the following parameters:
  - rentals?price_min
  - rentals?price_max
  - rentals?limit
  - rentals?ids
  - rentals?offset
  - rentals?near
  - rentals?sort
  - combinations of the above

## How to run the project locally

### Before you start:

You have to provide your own config.env in the main project directory file:
- APP_ENV=LOC/DEBUG/PROD
- DB_HOSTS
- DB_NAME
- DB_USERNAME
- DB_PASSWORD
- DB_PORT

### How to start the server

1. `docker-compose up` to start and populate the database
2. `go run main.go` to start the application
3. `go test -v ./...` to run the tests