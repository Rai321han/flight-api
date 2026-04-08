package validators

import (
	"flight-api/models"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var allowedSortFields = map[string]bool{
	"AvgTicketPrice": true,
	"timestamp":      true,
}

var allowedParams = map[string]bool{
	"carrier": true, "flightNum": true,
	"destCountry": true, "originCountry": true,
	"destCity": true, "originCity": true,
	"destAirportID": true, "originAirportID": true,
	"destAirport": true, "originAirport": true,
	"dateFrom": true, "dateTo": true,
	"priceMin": true, "priceMax": true,
	"cancelled": true,
	"destLat":   true, "destLon": true,
	"originLat": true, "originLon": true,
	"sortBy": true, "order": true,
	"limit": true, "page": true,
}

const (
	defaultSortBy = "timestamp"
	defaultOrder  = "asc"
	defaultLimit  = 10
	defaultPage   = 1
	maxLimit      = 1000
	dateLayout    = "2006-01-02"
)

func ParseAndValidateFlightFilters(
	getParam func(string) string,
	query url.Values,
) (*models.FlightFilters, []string) {
	errs := validateAllowedQueryParams(query)

	f := &models.FlightFilters{
		SortBy: defaultSortBy,
		Order:  defaultOrder,
		Limit:  defaultLimit,
		Page:   defaultPage,
	}

	applyStringFilters(f, getParam)
	parseDateFilters(f, getParam, &errs)
	parsePriceFilters(f, getParam, &errs)
	parsePaginationFilters(f, getParam, &errs)

	if len(errs) > 0 {
		return nil, errs
	}
	return f, nil
}

func validateAllowedQueryParams(query url.Values) []string {
	var errs []string
	for key := range query {
		if !allowedParams[key] {
			errs = append(errs, fmt.Sprintf("unsupported parameter: %s", key))
		}
	}
	return errs
}

func applyStringFilters(f *models.FlightFilters, getParam func(string) string) {
	f.Carrier = optionalString(getParam("carrier"))
	f.FlightNum = optionalString(getParam("flightNum"))
	f.DestCountry = optionalString(getParam("destCountry"))
	f.OriginCountry = optionalString(getParam("originCountry"))
	f.DestCity = optionalString(getParam("destCity"))
	f.OriginCity = optionalString(getParam("originCity"))
	f.DestAirportID = optionalString(getParam("destAirportID"))
	f.OriginAirportID = optionalString(getParam("originAirportID"))
	f.DestAirport = optionalString(getParam("destAirport"))
	f.OriginAirport = optionalString(getParam("originAirport"))
}

func parseDateFilters(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("dateFrom"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			*errs = append(*errs, "dateFrom: must be in YYYY-MM-DD format")
		} else {
			f.DateFrom = &t
		}
	}

	if v := getParam("dateTo"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			*errs = append(*errs, "dateTo: must be in YYYY-MM-DD format")
		} else {
			f.DateTo = &t
		}
	}

	if f.DateFrom != nil && f.DateTo != nil && f.DateTo.Before(*f.DateFrom) {
		*errs = append(*errs, "dateTo: must be on or after dateFrom")
	}
}

func parsePriceFilters(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	f.PriceMin, _ = parseOptionalFloat(getParam("priceMin"), "priceMin", 0, -1, errs)
	f.PriceMax, _ = parseOptionalFloat(getParam("priceMax"), "priceMax", 0, -1, errs)

	if f.PriceMin != nil && *f.PriceMin < 0 {
		*errs = append(*errs, "priceMin: must be >= 0")
		f.PriceMin = nil
	}
	if f.PriceMax != nil && *f.PriceMax < 0 {
		*errs = append(*errs, "priceMax: must be >= 0")
		f.PriceMax = nil
	}
	if f.PriceMin != nil && f.PriceMax != nil && *f.PriceMax < *f.PriceMin {
		*errs = append(*errs, "priceMax: must be >= priceMin")
	}
}

func parsePaginationFilters(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			*errs = append(*errs, "limit: must be a positive integer")
		} else if n > maxLimit {
			*errs = append(*errs, fmt.Sprintf("limit: must not exceed %d", maxLimit))
		} else {
			f.Limit = n
		}
	}

	if v := getParam("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			*errs = append(*errs, "page: must be a positive integer")
		} else {
			f.Page = n
		}
	}
}

func optionalString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

func parseOptionalFloat(raw, field string, minVal, maxVal float64, errs *[]string) (*float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		*errs = append(*errs, fmt.Sprintf("%s: must be a valid number", field))
		return nil, false
	}
	if maxVal > minVal && (v < minVal || v > maxVal) {
		*errs = append(*errs, fmt.Sprintf("%s: must be between %.2f and %.2f", field, minVal, maxVal))
		return nil, false
	}
	return &v, true
}
