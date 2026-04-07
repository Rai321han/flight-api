package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// FlightFilters holds all validated, type-coerced query parameters.
// Pointer fields mean "not provided" when nil.
type FlightFilters struct {
	Carrier         *string
	FlightNum       *string
	DestCountry     *string
	OriginCountry   *string
	DestCity        *string
	OriginCity      *string
	DateFrom        *time.Time
	DateTo          *time.Time
	PriceMin        *float64
	PriceMax        *float64
	Cancelled       *bool
	DestAirportID   *string
	OriginAirportID *string
	DestAirport     *string
	OriginAirport   *string
	DestLat         *float64
	DestLon         *float64
	OriginLat       *float64
	OriginLon       *float64
	SortBy          string
	Order           string
	Limit           int
}

const (
	defaultSortBy = "timestamp"
	defaultOrder  = "asc"
	defaultLimit  = 10
	maxLimit      = 1000
	dateLayout    = "2006-01-02"
)

// allowedSortFields restricts sortBy to known ES fields, preventing injection.
var allowedSortFields = map[string]bool{
	"AvgTicketPrice": true,
	"timestamp":      true,
}

type MockFlightService struct{}

func NewMockFlightService() *MockFlightService {
	return &MockFlightService{}
}

func (m *MockFlightService) SearchFlights(query *FlightFilters) ([]Flight, error) {
	// Fake data so controller can run
	return []Flight{
		{Name: "Kibana-Airlines FX123", Price: 450},
		{Name: "Elastic-Air EA456", Price: 380},
		{Name: "Logstash-Jet LJ789", Price: 520},
	}, nil
}

type Flight struct {
	Name  string
	Price float64
}

type SearchController struct {
	beego.Controller
}

// SearchFlights handles GET /api/v1/search requests.
// It parses and validates query parameters, then delegates to the FlightService to perform the search.
// Finally, it returns JSON responses with appropriate status codes or error messages.
func (c *SearchController) SearchFlights() {
	filters, validationErrs := ParseAndValidateFlightFilters(c.Ctx.Request.URL.Query().Get)
	if len(validationErrs) > 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]any{
			"error":   "invalid query parameters",
			"details": validationErrs,
		}
		c.ServeJSON()
		return
	}

	results, err := NewMockFlightService().SearchFlights(filters)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]any{
			"error": "internal server error",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]any{
		"data": results,
	}
	c.ServeJSON()
}

// ParseAndValidateFlightFilters is decoupled from beego's Controller so it is
// directly unit-testable without spinning up an HTTP server.
// getParam mirrors url.Values.Get — pass c.Ctx.Request.URL.Query().Get in production,
// or a map lookup func in tests.
func ParseAndValidateFlightFilters(getParam func(string) string) (*FlightFilters, []string) {
	var errs []string
	f := &FlightFilters{
		SortBy: defaultSortBy,
		Order:  defaultOrder,
		Limit:  defaultLimit,
	}

	// --- String fields (no format constraint, just trimming) ---
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

	// --- Date fields ---
	if v := getParam("dateFrom"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			errs = append(errs, "dateFrom: must be in YYYY-MM-DD format")
		} else {
			f.DateFrom = &t
		}
	}
	if v := getParam("dateTo"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			errs = append(errs, "dateTo: must be in YYYY-MM-DD format")
		} else {
			f.DateTo = &t
		}
	}
	if f.DateFrom != nil && f.DateTo != nil && f.DateTo.Before(*f.DateFrom) {
		errs = append(errs, "dateTo: must be on or after dateFrom")
	}

	// --- Price fields ---
	f.PriceMin, _ = parseOptionalFloat(getParam("priceMin"), "priceMin", 0, -1, &errs)
	f.PriceMax, _ = parseOptionalFloat(getParam("priceMax"), "priceMax", 0, -1, &errs)
	if f.PriceMin != nil && *f.PriceMin < 0 {
		errs = append(errs, "priceMin: must be >= 0")
		f.PriceMin = nil
	}
	if f.PriceMax != nil && *f.PriceMax < 0 {
		errs = append(errs, "priceMax: must be >= 0")
		f.PriceMax = nil
	}
	if f.PriceMin != nil && f.PriceMax != nil && *f.PriceMax < *f.PriceMin {
		errs = append(errs, "priceMax: must be >= priceMin")
	}

	// --- Bool field ---
	if v := getParam("cancelled"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			errs = append(errs, "cancelled: must be 'true' or 'false'")
		} else {
			f.Cancelled = &b
		}
	}

	// --- Geo fields: dest ---
	destLat, destLatOk := parseOptionalFloat(getParam("destLat"), "destLat", -90, 90, &errs)
	destLon, destLonOk := parseOptionalFloat(getParam("destLon"), "destLon", -180, 180, &errs)
	if destLatOk != destLonOk {
		errs = append(errs, "destLat and destLon must both be provided together or not at all")
	} else {
		f.DestLat = destLat
		f.DestLon = destLon
	}

	// --- Geo fields: origin ---
	originLat, originLatOk := parseOptionalFloat(getParam("originLat"), "originLat", -90, 90, &errs)
	originLon, originLonOk := parseOptionalFloat(getParam("originLon"), "originLon", -180, 180, &errs)
	if originLatOk != originLonOk {
		errs = append(errs, "originLat and originLon must both be provided together or not at all")
	} else {
		f.OriginLat = originLat
		f.OriginLon = originLon
	}

	// --- sortBy ---
	if v := getParam("sortBy"); v != "" {
		if !allowedSortFields[v] {
			allowed := make([]string, 0, len(allowedSortFields))
			for k := range allowedSortFields {
				allowed = append(allowed, k)
			}
			errs = append(errs, fmt.Sprintf("sortBy: must be one of [%s]", strings.Join(allowed, ", ")))
		} else {
			f.SortBy = v
		}
	}

	// --- order ---
	if v := getParam("order"); v != "" {
		v = strings.ToLower(v)
		if v != "asc" && v != "desc" {
			errs = append(errs, "order: must be 'asc' or 'desc'")
		} else {
			f.Order = v
		}
	}

	// --- limit ---
	if v := getParam("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			errs = append(errs, "limit: must be a positive integer")
		} else if n > maxLimit {
			errs = append(errs, fmt.Sprintf("limit: must not exceed %d", maxLimit))
		} else {
			f.Limit = n
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return f, nil
}

// --- helpers ---

func optionalString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

// parseOptionalFloat parses a float from a raw string and validates it is within
// [min, max] when max > min. Returns (value, wasProvided).
// minVal/maxVal: pass equal values (e.g. 0,0 or -1,-1) to skip range check.
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
