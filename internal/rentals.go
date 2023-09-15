package internal

import (
	"outdoorsy-api/database"
)

// GetASingleRental retrieves a single rental from the database by the specified id.
func GetASingleRental(id int) (rental Rental, err error) {
	err = database.GetSingleRecordNamedQuery(&rental, selectSingleRentalQuery, map[string]interface{}{"id": id})
	return
}
