package internal

import (
	"net/url"
	"outdoorsy-api/database"
)

// GetASingleRental retrieves a single rental from the database by the specified id.
func GetASingleRental(id int) (rental Rental, err error) {
	err = database.GetSingleRecordNamedQuery(&rental, selectSingleRentalQuery, map[string]interface{}{"id": id})
	return
}

func GetMultipleRentals(params url.Values) (rentals []Rental, err error) {
	if len(params) == 0 {
		err = database.GetMultipleRecords(&rentals, selectAllRentalsQuery)
		return
	}
	return
}
