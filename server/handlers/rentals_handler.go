package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"outdoorsy-api/internal"
	"outdoorsy-api/utils"
	"strconv"
)

func SingleRentalHandler(ginCtx *gin.Context) {
	idAsString, _ := ginCtx.Params.Get("id")

	id, err := strconv.Atoi(idAsString)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]string{"error": "provided parameter should be of type int"})
		return
	}

	rental, err := internal.GetASingleRental(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusNoContent, rental)
			return
		}
		utils.GetLogger().WithFields(log.Fields{"error": err.Error(), "id": id}).Error("Error on getting single rental from the database")
		ginCtx.JSON(http.StatusInternalServerError, rental)
		return
	}

	ginCtx.JSON(http.StatusOK, rental)
}

func MultipleRentalsHandler(ginCtx *gin.Context) {
	rental, err := internal.GetMultipleRentals(ginCtx.Request.URL.Query())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusNoContent, rental)
			return
		}

		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on getting all rentals from the database")
		ginCtx.JSON(http.StatusInternalServerError, rental)
		return
	}

	ginCtx.JSON(http.StatusOK, rental)
}
