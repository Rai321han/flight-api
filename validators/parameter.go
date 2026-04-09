package validators

import (
	"flight-api/models"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// allowedSortFields defines which fields can be used for sorting results.
var allowedSortFields = map[string]bool{
	"AvgTicketPrice": true,
	"timestamp":      true,
}

// allowedParams defines which query parameters are accepted by the search endpoint. Any parameter not in this list will cause a validation error.
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

// ParseAndValidateFlightFilters takes the raw query parameters from the request, validates them, and constructs a FlightFilters struct that can be used by the FlightService.
// It returns any validation errors encountered during parsing.
func ParseAndValidateFlightFilters(
	query url.Values,
) (*models.FlightFilters, []string) {
	errs := validateAllowedQueryParams(query)

	f := &models.FlightFilters{
		SortBy: defaultSortBy,
		Order:  defaultOrder,
		Limit:  defaultLimit,
		Page:   defaultPage,
	}

	applyStringFilters(f, query.Get)
	f.Cancelled = parseBoolParam("cancelled", query.Get, &errs)
	parseDateFilters(f, query.Get, &errs)
	parsePriceFilters(f, query.Get, &errs)
	parsePaginationFilters(f, query.Get, &errs)
	parseSortAndOrder(f, query.Get, &errs)
	parseGeoFilters(f, query.Get, &errs)

	if len(errs) > 0 {
		return nil, errs
	}
	return f, nil
}

// validateAllowedQueryParams checks if any query parameters are present that are not in the allowedParams list.
// It returns a slice of error messages for any unsupported parameters found.
func validateAllowedQueryParams(query url.Values) []string {
	var errs []string
	for key := range query {
		if !allowedParams[key] {
			errs = append(errs, fmt.Sprintf("unsupported parameter: %s", key))
		}
	}
	return errs
}

// applyStringFilters extracts string-based filters from the query parameters and assigns them to the FlightFilters struct.
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

// parseBoolParam attempts to parse a boolean query parameter. It returns a pointer to the boolean value if successful, or nil if the parameter is not provided.
func parseBoolParam(param string, getParam func(string) string, errs *[]string) *bool {
	value := strings.TrimSpace(getParam(param))
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return nil
	}
	if value == "true" {
		v := true
		return &v
	}
	if value == "false" {
		v := false
		return &v
	}
	*errs = append(*errs, fmt.Sprintf("%s must be true or false", param))
	return nil
}

func parseGeoFilters(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	f.DestLat, _ = parseOptionalFloat(getParam("destLat"), "destLat", -90, 90, errs)
	f.DestLon, _ = parseOptionalFloat(getParam("destLon"), "destLon", -180, 180, errs)
	f.OriginLat, _ = parseOptionalFloat(getParam("originLat"), "originLat", -90, 90, errs)
	f.OriginLon, _ = parseOptionalFloat(getParam("originLon"), "originLon", -180, 180, errs)
}

func parseSortAndOrder(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("sortBy"); v != "" {
		if !allowedSortFields[v] {
			*errs = append(*errs, fmt.Sprintf("sortBy: unsupported sort field '%s'", v))
		} else {
			f.SortBy = v
		}
	}

	if v := getParam("order"); v != "" {
		v = strings.ToLower(v)
		if v != "asc" && v != "desc" {
			*errs = append(*errs, "order: must be 'asc' or 'desc'")
		} else {
			f.Order = v
		}
	}
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
		*errs = append(*errs, "priceMin: must be greater than 0")
		f.PriceMin = nil
	}
	if f.PriceMax != nil && *f.PriceMax < 0 {
		*errs = append(*errs, "priceMax: must be greater than 0")
		f.PriceMax = nil
	}
	if f.PriceMin != nil && f.PriceMax != nil && *f.PriceMax < *f.PriceMin {
		*errs = append(*errs, "priceMax: must be greater than priceMin")
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
