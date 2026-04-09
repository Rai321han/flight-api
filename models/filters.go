package models

import "time"

// FlightFilters holds all optional search criteria, sort preferences, and
// pagination settings parsed from an incoming HTTP request. Pointer fields
// are nil when the caller did not supply that parameter, allowing the query
// builder to skip the corresponding Elasticsearch filter clause entirely.
type FlightFilters struct {
	// String filters — exact match against the flights index.
	// Pointer is nil when the parameter was absent from the request.
	Carrier         *string // airline carrier name
	FlightNum       *string // flight number identifier
	DestCountry     *string // destination country code (e.g. "US", "DE")
	OriginCountry   *string // origin country code
	DestCity        *string // destination city name
	OriginCity      *string // origin city name
	DestAirportID   *string // destination airport IATA code (e.g. "JFK")
	OriginAirportID *string // origin airport IATA code
	DestAirport     *string // destination airport full name
	OriginAirport   *string // origin airport full name

	// Date range filters applied to the timestamp field.
	// Either bound may be set independently; both nil means no date filter.
	DateFrom *time.Time // lower bound, inclusive (gte), format YYYY-MM-DD
	DateTo   *time.Time // upper bound, inclusive (lte), format YYYY-MM-DD

	// Price range filters applied to the AvgTicketPrice field.
	// Either bound may be set independently; both nil means no price filter.
	PriceMin *float64 // lower bound, inclusive (gte), must be >= 0
	PriceMax *float64 // upper bound, inclusive (lte), must be >= PriceMin

	// Boolean filter on the Cancelled field.
	// true returns only cancelled flights; false returns only active flights.
	// nil means no cancelled filter is applied.
	Cancelled *bool

	// Geo-distance filters — each pair restricts results to documents whose
	// corresponding geo-point field lies within 100 km of the coordinate.
	// Both lat and lon must be non-nil for the filter to be applied.
	DestLat   *float64 // destination latitude  (-90  to  90), requires DestLon
	DestLon   *float64 // destination longitude (-180 to 180), requires DestLat
	OriginLat *float64 // origin latitude       (-90  to  90), requires OriginLon
	OriginLon *float64 // origin longitude      (-180 to 180), requires OriginLat

	// Sort and pagination — always populated by the validator with defaults
	// when the caller omits them; never nil or zero after successful parsing.
	SortBy string // sort field: "timestamp" (default) or "AvgTicketPrice"
	Order  string // sort direction: "asc" (default) or "desc"
	Limit  int    // maximum documents per page (1..1000, default 10)
	Page   int    // 1-based page number (default 1)
}