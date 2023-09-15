package internal

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"net/url"
	"outdoorsy-api/database"
	"testing"
)

type configurations struct {
	DBHosts    string `json:"db_hosts" koanf:"DB_HOSTS" valid:"required"`
	DBUsername string `json:"db_username" koanf:"DB_USERNAME" valid:"required"`
	DBPassword string `json:"db_password" koanf:"DB_PASSWORD" valid:"required"`
	DBPort     string `json:"db_port" koanf:"DB_PORT" valid:"required"`
	DBName     string `json:"db_name" koanf:"DB_NAME" valid:"required"`
}

func setupTest(test *testing.T) func() {
	var (
		parser = koanf.New(".")
		config configurations
	)

	err := parser.Load(file.Provider("../config.env"), dotenv.Parser())
	if err != nil {
		test.Fatal(err.Error())
	}

	err = parser.Unmarshal("", &config)
	if err != nil {
		test.Fatal(err.Error())
	}

	_, err = validator.ValidateStruct(config)
	if err != nil {
		test.Fatal(err.Error())
	}

	database.Init(
		config.DBHosts,
		config.DBUsername,
		config.DBPassword,
		config.DBPort,
		config.DBName,
	)

	return func() {
	}
}

func TestGetASingleRentalShouldReturnCorrectResponse(test *testing.T) {
	defer setupTest(test)()
	rentalId := 1
	rentalName := "'Abaco' VW Bay Window: Westfalia Pop-top"

	rental, err := GetASingleRental(rentalId)
	if err != nil {
		test.Fatalf("Error on getting rental with id %d", rentalId)
	}

	assert.Equal(test, rentalId, rental.IdRental, "Expected rental id and request id to match")
	assert.Equal(test, rentalName, rental.Name, "Expected rental name to match the retrieved record")
}

func TestGetASingleRentalShouldReturnNoRowsInResultSetIfIdIsInvalid(test *testing.T) {
	defer setupTest(test)()
	rentalId := -1

	rental, err := GetASingleRental(rentalId)
	if err == nil {
		test.Fatalf("Error should be returned on getting rental with invalid id - %d", rentalId)
	}

	assert.Error(test, err, "Error should be returned on getting rental with invalid id.")
	assert.Equal(test, err.Error(), "sql: no rows in result set", "Correct error message should be returned if result was not found")
	assert.Equal(test, rental, Rental{}, "Empty rental should be returned as result from select with invalid id")
}

func TestGetMultipleRentalsShouldReturnAllRentalsIfNoParametersArePassed(test *testing.T) {
	defer setupTest(test)()

	var (
		params               url.Values
		expectedRentalsCount = 30
	)

	rentals, failedValidation, err := GetMultipleRentals(params)
	if err != nil {
		test.Fatal("Getting all rentals with no parameters passed should return correct result")
	}

	if failedValidation {
		test.Fatal("Getting all rentals with no parameters passed should not fail validations")
	}

	assert.Equal(test, expectedRentalsCount, len(rentals), "There should be 30 default rentals in the database")
}

func TestGetMultipleRentalsShouldReturnDescriptiveErrorInCaseOfFailedMinPriceParameterValidation(test *testing.T) {
	defer setupTest(test)()

	var (
		params                           = make(url.Values)
		expectedFailedPriceValidationMsg = "price must be a positive number"
	)

	params.Set("price_min", "-1")

	rentals, failedValidation, err := GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid price should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid price should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedPriceValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")
}

func TestGetMultipleRentalsShouldReturnDescriptiveErrorInCaseOfFailedMaxPriceParameterValidation(test *testing.T) {
	defer setupTest(test)()

	var (
		params                           = make(url.Values)
		expectedFailedPriceValidationMsg = "price must be a positive number"
	)

	params.Set("price_max", "-1")

	rentals, failedValidation, err := GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid price should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid price should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedPriceValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")
}

func TestGetMultipleRentalsShouldReturnDescriptiveErrorInCaseOfBiggerOrEqualMinPriceThanMaxPrice(test *testing.T) {
	defer setupTest(test)()

	var (
		params                           = make(url.Values)
		expectedFailedPriceValidationMsg = "price_min must be less than price_max"
	)

	params.Set("price_max", "1000")
	params.Set("price_min", "1000.01")

	rentals, failedValidation, err := GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with bigger or equal min price to max price should return error")
	assert.True(test, failedValidation, "Failed validation is expected")
	assert.Equal(test, err.Error(), expectedFailedPriceValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("price_min", "1000")
	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with bigger or equal min price to max price should return error")
	assert.True(test, failedValidation, "Failed validation is expected")
	assert.Equal(test, err.Error(), expectedFailedPriceValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")
}

func TestGetMultipleRentalsShouldReturnDescriptiveErrorInCaseOfFailedLimitParameterValue(test *testing.T) {
	defer setupTest(test)()

	var (
		params                      = make(url.Values)
		expectedFailedValidationMsg = "limit and offset must be a positive integer"
	)

	params.Set("limit", "0")

	rentals, failedValidation, err := GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid limit should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid limit should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("limit", "1.5")

	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid limit should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid limit should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("limit", "-1")

	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid limit should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid limit should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsg, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")
}

func TestGetMultipleRentalsShouldReturnDescriptiveErrorInCaseOfFailedIdParameterValidation(test *testing.T) {
	defer setupTest(test)()

	var (
		params                                       = make(url.Values)
		expectedFailedValidationMsgForNonInt         = "ids should be int values"
		expectedFailedValidationMsgForNonPositiveInt = "ids should be positive numbers greater than 0"
	)

	params.Set("ids", "1.2,2,5")

	rentals, failedValidation, err := GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid ids should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid ids should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsgForNonInt, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("ids", "holla,cute code")

	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid ids should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid ids should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsgForNonInt, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("ids", "1,10,-1")

	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid ids should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid ids should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsgForNonPositiveInt, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")

	params.Set("ids", "0")

	rentals, failedValidation, err = GetMultipleRentals(params)

	assert.Error(test, err, "Getting rentals with invalid ids should return error")
	assert.True(test, failedValidation, "Getting rentals with invalid ids should return true for failed validation")
	assert.Equal(test, err.Error(), expectedFailedValidationMsgForNonPositiveInt, "Correct error message is expected in case of failed validation")
	assert.Equal(test, 0, len(rentals), "No results should be returned in case of failed validation")
}

func TestGetMultipleRentalsWithIdParameterShouldReturnCorrectResponse(test *testing.T) {
	defer setupTest(test)()

	var (
		params                = make(url.Values)
		expectedRentalResults = 3
	)
	params.Set("ids", "1,6,11")

	rentals, failedValidation, err := GetMultipleRentals(params)
	if err != nil {
		test.Fatalf("Error on getting rentals with ids - %s", params.Get("ids"))
	}

	if failedValidation {
		test.Fatalf("There should be no failed validation in case of valid input")
	}

	assert.Equal(test, expectedRentalResults, len(rentals), "Expected 3 results to be retrieved")
	assert.Equal(test, 1, rentals[0].IdRental, "Expected retrieved rentals ids to match")
	assert.Equal(test, 6, rentals[1].IdRental, "Expected retrieved rentals ids to match")
	assert.Equal(test, 11, rentals[2].IdRental, "Expected retrieved rentals ids to match")
}

func TestGetMultipleRentalsWithSortParameterShouldSortByPrice(test *testing.T) {
	defer setupTest(test)()

	var params = make(url.Values)
	params.Set("sort", "price")

	rentals, failedValidation, err := GetMultipleRentals(params)
	if err != nil {
		test.Fatalf("Error on getting rentals with sorting parameter! - %s", err.Error())
	}

	if failedValidation {
		test.Fatalf("There should be no failed validation in case of a valid input")
	}

	assert.Equal(test, 12, rentals[0].IdRental, "Expected retrieved rentals id 12 to be first result")
	assert.Equal(test, 9, rentals[1].IdRental, "Expected retrieved rentals id 9 to be first result")
	assert.Equal(test, 13, rentals[2].IdRental, "Expected retrieved rentals id 13 to be first result")
}

func TestGetMultipleRentalsWithPriceMinShouldNotReturnCheaperVehicles(test *testing.T) {
	defer setupTest(test)()

	var params = make(url.Values)
	params.Set("price_min", "10000")
	params.Set("sort", "price")

	rentals, failedValidation, err := GetMultipleRentals(params)
	if err != nil {
		test.Fatalf("Error on getting rentals with sorting parameter! - %s", err.Error())
	}

	if failedValidation {
		test.Fatalf("There should be no failed validation in case of a valid input")
	}

	expectedResult := rentals[0].Price.Day >= 10000

	assert.True(test, expectedResult, "Expected the first result to be with min price of 10k")

	params.Set("price_min", "5000")

	rentals, failedValidation, err = GetMultipleRentals(params)

	expectedResult = rentals[0].Price.Day >= 5000

	assert.True(test, expectedResult, "Expected the first result to be with min price of 5k")
}
