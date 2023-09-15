package server

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"outdoorsy-api/database"
	"outdoorsy-api/server/handlers"
	"testing"
)

type configurations struct {
	DBHosts    string `json:"db_hosts" koanf:"DB_HOSTS" valid:"required"`
	DBUsername string `json:"db_username" koanf:"DB_USERNAME" valid:"required"`
	DBPassword string `json:"db_password" koanf:"DB_PASSWORD" valid:"required"`
	DBPort     string `json:"db_port" koanf:"DB_PORT" valid:"required"`
	DBName     string `json:"db_name" koanf:"DB_NAME" valid:"required"`
}

var router *gin.Engine

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

	router = gin.New()
	router.GET("/healths", handlers.HealthCheck)
	router.GET("/metrics", handlers.Metrics)
	router.GET("/rentals/:id", handlers.SingleRentalHandler)
	router.GET("/rentals", handlers.MultipleRentalsHandler)

	return func() {
	}
}

func TestSingleRentalHandler(test *testing.T) {
	defer setupTest(test)()

	test.Run("SuccessfulRequest", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/rentals/1", nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		responseRecorder := httptest.NewRecorder()

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
	})

	test.Run("RentalNotFound", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/rentals/999", nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		responseRecorder := httptest.NewRecorder()

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
	})

	test.Run("BadRentalRequest", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/rentals/0.5", nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		responseRecorder := httptest.NewRecorder()

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	})
}
