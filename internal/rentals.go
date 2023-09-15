package internal

import (
	"fmt"
	"net/url"
	"outdoorsy-api/database"
	"outdoorsy-api/utils"
	"strings"
)

// GetASingleRental retrieves a single rental from the database by the specified id.
func GetASingleRental(id int) (rental Rental, err error) {
	err = database.GetSingleRecordNamedQuery(&rental, selectSingleRentalQuery, map[string]interface{}{"id": id})
	return
}

// GetMultipleRentals will check for URL parameters. If there are no parameters - all records from the database
// will be retrieved. If parameters are provided - they will be validated and if valid - filtration will be made
// If the parameters are not valid - failedValidation will be set to true and descriptive validation error will
// be returned. In case of non supported parameter - all records will be retrieved as if no parameter was
// added.
func GetMultipleRentals(params url.Values) (rentals []Rental, failedValidation bool, err error) {
	if len(params) == 0 {
		err = database.GetMultipleRecords(&rentals, selectAllRentalsQuery)
		return
	}

	if err = utils.ValidateParameters(params); err != nil {
		failedValidation = true
		return
	}

	additionalQueryParams := transpileParamsToDBQueries(params)
	err = database.GetMultipleRecords(&rentals, selectAllRentalsQuery+additionalQueryParams)
	return
}

func transpileParamsToDBQueries(params url.Values) (additionalQueryParams string) {
	var builder strings.Builder

	handleDBQueryWhereClauses(params, &builder)
	handleDBQuerySortingClause(params, &builder)
	handleDBQueryFinalClauses(params, &builder)

	additionalQueryParams = builder.String()
	additionalQueryParams = strings.ReplaceAll(additionalQueryParams, "[", "")
	additionalQueryParams = strings.ReplaceAll(additionalQueryParams, "]", "")
	return
}

func handleDBQueryWhereClauses(params url.Values, builder *strings.Builder) {
	var (
		queryWhereClause string
		whereParameters  = map[string]string{
			"price_min": " price_per_day >= %s",
			"price_max": " price_per_day <= %s",
			"ids":       " rentals.id IN (%s)",
			"near":      " lat >= %s AND lng >= %s",
		}
	)

	for key, value := range params {
		content, contentWasFound := whereParameters[key]
		if !contentWasFound || len(value) > 0 && value[0] == "" {
			continue
		}

		if key == "near" {
			splitValues := strings.Split(value[0], ",")
			queryWhereClause = fmt.Sprintf(content, splitValues[0], splitValues[1])
		} else {
			queryWhereClause = fmt.Sprintf(content, value)
		}

		if builder.Len() > 0 {
			builder.WriteString(" AND ")
			builder.WriteString(queryWhereClause)
		} else {
			builder.WriteString(" WHERE ")
			builder.WriteString(queryWhereClause)
		}
	}
}

func handleDBQuerySortingClause(params url.Values, builder *strings.Builder) {
	var querySortingClause string

	if params.Has("sort") {
		sortingValue := params.Get("sort")
		if sortingValue == "" {
			return
		} else if sortingValue == "price" {
			querySortingClause = " ORDER BY price_per_day"
		} else {
			querySortingClause = fmt.Sprintf(" ORDER BY %s", sortingValue)
		}
		builder.WriteString(querySortingClause)
	}
}

func handleDBQueryFinalClauses(params url.Values, builder *strings.Builder) {
	var (
		queryFinalClause      string
		finalClauseParameters = map[string]string{
			"limit":  " LIMIT %s",
			"offset": " OFFSET %s",
		}
	)

	for key, value := range params {
		content, contentWasFound := finalClauseParameters[key]
		if !contentWasFound || len(value) > 0 && value[0] == "" {
			continue
		}
		queryFinalClause = fmt.Sprintf(content, value)
		builder.WriteString(queryFinalClause)
	}
}
