package validators

import (
	"flight-api/models"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// allowedSortFields restricts sortBy to known ES fields, preventing injection.
var allowedSortFields = map[string]bool{
	"AvgTicketPrice": true,
	"timestamp":      true,
}

const (
	defaultSortBy = "timestamp"
	defaultOrder  = "asc"
	defaultLimit  = 10
	maxLimit      = 1000
	dateLayout    = "2006-01-02"
)

// ParseAndValidateFlightFilters is decoupled from beego's Controller so it is
// directly unit-testable without spinning up an HTTP server.
// getParam mirrors url.Values.Get — pass c.Ctx.Request.URL.Query().Get in production,
// or a map lookup func in tests.
// ParseAndValidateFlightFilters parses and validates all query parameters.
func ParseAndValidateFlightFilters(getParam func(string) string) (*models.FlightFilters, error) {
	var errs []string
	f := &models.FlightFilters{
		SortBy: defaultSortBy,
		Order:  defaultOrder,
		Limit:  defaultLimit,
	}

	// Parse all field categories
	parseStringFields(f, getParam)
	parseAndValidateDateFields(f, getParam, &errs)
	parseAndValidatePriceFields(f, getParam, &errs)
	parseAndValidateBoolField(f, getParam, &errs)
	parseAndValidateGeoFields(f, getParam, &errs)
	parseAndValidateSortBy(f, getParam, &errs)
	parseAndValidateOrder(f, getParam, &errs)
	parseAndValidateLimit(f, getParam, &errs)

	if len(errs) > 0 {
		msg := "invalid query parameters"
		if len(errs) == 1 {
			msg = errs[0]
		}
		return nil, models.NewAppError(400, msg, nil)
	}
	return f, nil
}

// parseStringFields sets optional string query parameters.
func parseStringFields(f *models.FlightFilters, getParam func(string) string) {
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

// parseAndValidateDateFields validates date range parameters.
func parseAndValidateDateFields(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
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

// parseAndValidatePriceFields validates price range parameters.
func parseAndValidatePriceFields(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
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

// parseAndValidateBoolField validates the cancelled boolean parameter.
func parseAndValidateBoolField(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("cancelled"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			*errs = append(*errs, "cancelled: must be 'true' or 'false'")
		} else {
			f.Cancelled = &b
		}
	}
}

// parseAndValidateGeoFields validates geo location parameters.
func parseAndValidateGeoFields(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	destLat, destLatOk := parseOptionalFloat(getParam("destLat"), "destLat", -90, 90, errs)
	destLon, destLonOk := parseOptionalFloat(getParam("destLon"), "destLon", -180, 180, errs)
	if destLatOk != destLonOk {
		*errs = append(*errs, "destLat and destLon must both be provided together or not at all")
	} else if destLatOk {
		f.DestLat = destLat
		f.DestLon = destLon
	}

	originLat, originLatOk := parseOptionalFloat(getParam("originLat"), "originLat", -90, 90, errs)
	originLon, originLonOk := parseOptionalFloat(getParam("originLon"), "originLon", -180, 180, errs)
	if originLatOk != originLonOk {
		*errs = append(*errs, "originLat and originLon must both be provided together or not at all")
	} else if originLatOk {
		f.OriginLat = originLat
		f.OriginLon = originLon
	}
}

// parseAndValidateSortBy validates the sortBy parameter.
func parseAndValidateSortBy(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("sortBy"); v != "" {
		if allowedSortFields[v] {
			f.SortBy = v
		} else {
			allowed := make([]string, 0, len(allowedSortFields))
			for k := range allowedSortFields {
				allowed = append(allowed, k)
			}
			*errs = append(*errs, fmt.Sprintf("sortBy: must be one of [%s]", strings.Join(allowed, ", ")))
		}
	}
}

// parseAndValidateOrder validates the order parameter.
func parseAndValidateOrder(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
	if v := getParam("order"); v != "" {
		v = strings.ToLower(v)
		if v == "asc" || v == "desc" {
			f.Order = v
		} else {
			*errs = append(*errs, "order: must be 'asc' or 'desc'")
		}
	}
}

// parseAndValidateLimit validates the limit parameter.
func parseAndValidateLimit(f *models.FlightFilters, getParam func(string) string, errs *[]string) {
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
}

// optionalString trims whitespace and returns a pointer to the string, or nil if empty.
func optionalString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

// parseOptionalFloat parses a float query parameter, validates it against optional min/max bounds,
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
