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

The biggest challenge here is to write "production ready code" and i won't consider my code production ready yet, because it's something that will take more than 4-5 hours to implement :)

The meaning of production ready code for me involves those main requirements:

1. Application stability - there should be no panics in case of errors or lost connections. In every case - the application should remain stable and handle requests. Fallbacks in case the database is down and reconnect mechanisms can help with this :)
2. Validations - validating incoming data can have many benefits like reduced outer systems load (like databases) - in case of invalid input we should not even consider using the db. And what about SQL injections? :)
3. Testing - Unit and integration tests, load tests, stress tests (shutting down external components and testing the reconnect mechs). We can verify the business logic, the reconnect mechanisms, the fallback mechanisms and many more thanks to a good test coverage and structure.
4. Monitoring - things like metrics, healths, logs should be implemented in order to ensure application stability and low running costs in the long run.
5. System security - excuse me for the * origin part :D DDOS attacks protection, access rules, SQL injections, CSRF attacks, hashing .. there are so many things to take into consideration here.
6. And finally after we ensure the above points i like to think about performance optimizations and performance monitoring. There is always something to improve and it's critical for modern day applications to be performant. Slow responses can affect the conversion rate of the website, the SEO and can have significant impact on user experience.

Keeping all of these points in mind if i had more time i would implement:

* More validations (there are couple more edge cases that can be considered)
* More tests (both unit and integration)
* Fallback responses in case of error or timeout of the database
* Load tests (using locust for example)
* Security - it's nice to have a devops / network security team but since i don't have it on my side i would have to do the things on my own. I would definitely take some courses in this topic and start to implement some goodies. :)
* Caching - most of the data seems to be with static nature, we can decrease the DB load and overall performance by using something like Redis

## REST API endpoints

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
