package utils

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func ValidateParameters(params url.Values) (err error) {
	var (
		priceMin = params.Get("price_min")
		priceMax = params.Get("price_max")
		limit    = params.Get("limit")
		offset   = params.Get("offset")
		ids      = params.Get("ids")
		near     = params.Get("near")
		minPrice float64
		maxPrice float64
	)

	if priceMin != "" {
		minPrice, err = validatePrice(priceMin)
		if err != nil {
			return
		}
	}
	if priceMax != "" {
		maxPrice, err = validatePrice(priceMax)
		if err != nil {
			return
		}
	}
	if priceMin != "" && priceMax != "" {
		err = validateMinAndMaxPrice(minPrice, maxPrice)
		if err != nil {
			return
		}
	}
	if limit != "" {
		err = validateIntegerValues(limit)
		if err != nil {
			return
		}
	}
	if offset != "" {
		err = validateIntegerValues(offset)
		if err != nil {
			return
		}
	}
	if ids != "" {
		err = validateArray(ids)
		if err != nil {
			return
		}
	}
	if near != "" {
		err = validateNear(near)
		if err != nil {
			return
		}
	}
	return
}

func validatePrice(price string) (priceAsNumber float64, err error) {
	priceAsNumber, err = strconv.ParseFloat(price, 64)
	if err != nil || priceAsNumber < 0 {
		err = errors.New("price must be a positive number")
		return
	}
	return
}

func validateMinAndMaxPrice(min float64, max float64) error {
	if min >= max {
		return errors.New("price_min must be less than price_max")
	} else if max <= min {
		return errors.New("price_max must be more than price_min")
	}
	return nil
}

func validateIntegerValues(limit string) error {
	num, err := strconv.Atoi(limit)
	if err != nil || num <= 0 {
		return errors.New("limit and offset must be a positive integer")
	}
	return nil
}

func validateArray(ids string) error {
	stringArr := strings.Split(ids, ",")
	for _, value := range stringArr {
		num, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("ids should be int values")
		}
		if num <= 0 {
			return errors.New("ids should be positive numbers greater than 0")
		}
	}
	return nil
}

func validateNear(nearContent string) error {
	if !strings.Contains(nearContent, ",") {
		return errors.New("there should be comma separator for the near parameter")
	}

	nearValues := strings.Split(nearContent, ",")
	if len(nearValues) != 2 {
		return errors.New("near values should be a comma separated pair")
	}

	for _, value := range nearValues {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("near values should be numeric")
		}
	}

	return nil
}
