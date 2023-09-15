package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthResponse struct {
	AppStatus      string `json:"AppStatus"`
	DatabaseStatus string `json:"Database"`
}

var health healthResponse

func HealthCheck(ginCtx *gin.Context) {
	checkDB(&health)
	ginCtx.JSON(http.StatusOK, health)
}

func checkDB(response *healthResponse) {
	response.AppStatus = "Healthy"
	response.DatabaseStatus = "Healthy"
}
