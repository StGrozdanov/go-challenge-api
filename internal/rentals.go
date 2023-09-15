package internal

import (
	"fmt"
	"net/url"
	"outdoorsy-api/database"
	"strings"
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
		if !contentWasFound {
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
		get := params.Get("sort")
		if get == "price" {
			querySortingClause = " ORDER BY price_per_day"
		} else {
			querySortingClause = fmt.Sprintf(" ORDER BY %s", get)
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
		if !contentWasFound {
			continue
		}
		queryFinalClause = fmt.Sprintf(content, value)
		builder.WriteString(queryFinalClause)
	}
}
