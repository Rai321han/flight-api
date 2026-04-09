package services

import "flight-api/models"

const flightIndex = "kibana_sample_data_flights"

// flightSourceFields defines the explicit _source projection included in every
// Elasticsearch search request. Declaring the projection here reduces network
// payload by omitting unmapped or internal fields, and serves as the single
// authoritative list of fields the search endpoint exposes to callers.
// Extend or trim this slice to control the response shape without touching
// query construction logic.
var flightSourceFields = []string{
	"FlightNum", "Carrier", "FlightDelay", "FlightDelayMin", "FlightDelayType",
	"OriginCityName", "OriginCountry", "OriginAirportID", "OriginAirport",
	"OriginLocation",
	"DestCityName", "DestCountry", "DestAirportID", "DestAirport",
	"DestLocation",
	"AvgTicketPrice", "Cancelled", "timestamp",
}

// termField maps a FlightFilters field to its corresponding Elasticsearch
// field name. The keyword flag controls whether the query targets the field
// directly (native keyword mapping) or via its .keyword sub-field (analysed
// text mapping). Consult the index mapping before changing this flag —
// appending .keyword to a native keyword field produces zero hits silently.
type termField struct {
	esField string
	value   *string
	keyword bool
}

// BuildSearchRequest constructs a complete Elasticsearch Search API request
// body from a validated FlightFilters value.
//
// Parameters:
//   - f: validated FlightFilters containing all active search criteria,
//     sort preference, and pagination settings
//
// Returns:
//   - map[string]any: complete ES request body ready for JSON serialisation,
//     including query, sort, from, size, and _source projection
//
// Defaults:
//   - sortBy defaults to "timestamp" when empty
//   - order  defaults to "asc" when empty
//   - from   is clamped to 0 to guard against negative offsets
func BuildSearchRequest(f *models.FlightFilters) map[string]any {
	filter := buildFilter(f)

	query := map[string]any{"match_all": map[string]any{}}
	if len(filter) > 0 {
		query = map[string]any{"bool": map[string]any{"filter": filter}}
	}

	sortField := f.SortBy
	if sortField == "" {
		sortField = "timestamp"
	}
	order := f.Order
	if order == "" {
		order = "asc"
	}

	from := (f.Page - 1) * f.Limit
	if from < 0 {
		from = 0
	}

	return map[string]any{
		"query":   query,
		"sort":    []map[string]any{{sortField: map[string]any{"order": order}}},
		"size":    f.Limit,
		"from":    from,
		"_source": flightSourceFields,
	}
}

// IndexName returns the Elasticsearch index name used for all flight queries.
//
// Returns:
//   - string: the index name constant for kibana_sample_data_flights
func IndexName() string {
	return flightIndex
}

// buildFilter assembles the bool.filter clause by delegating each filter
// concern to a dedicated builder function. Returns nil when no filters are
// active, which causes BuildSearchRequest to emit a match_all query instead.
//
// Filter groups:
//   - Term filters  : exact-match on string and ID fields
//   - Cancelled     : boolean flag filter
//   - Date range    : timestamp bounded by dateFrom and/or dateTo
//   - Price range   : AvgTicketPrice bounded by priceMin and/or priceMax
//   - Geo distance  : 100 km radius around destination and/or origin coordinates
func buildFilter(f *models.FlightFilters) []map[string]any {
	var filter []map[string]any
	filter = append(filter, buildTermFilters(f)...)
	filter = append(filter, buildCancelledFilter(f)...)
	filter = append(filter, buildDateRangeFilter(f)...)
	filter = append(filter, buildPriceRangeFilter(f)...)
	filter = append(filter, buildGeoFilter("DestLocation", f.DestLat, f.DestLon)...)
	filter = append(filter, buildGeoFilter("OriginLocation", f.OriginLat, f.OriginLon)...)
	return filter
}

// buildTermFilters returns one term clause per non-nil string filter field.
// Fields with native keyword mappings in the index are queried directly;
// fields with analysed text mappings are queried via their .keyword sub-field
// to enforce exact, case-sensitive matching.
//
// Parameters:
//   - carrier, flightNum, destCountry, originCountry  : native keyword mapping
//   - destAirportID, originAirportID                  : native keyword mapping
//   - destCity, originCity, destAirport, originAirport : text + .keyword sub-field
//
// Returns:
//   - []map[string]any: term clauses for all non-nil fields, or nil when none are set
func buildTermFilters(f *models.FlightFilters) []map[string]any {
	fields := []termField{
		{"Carrier", f.Carrier, false},
		{"FlightNum", f.FlightNum, false},
		{"DestCountry", f.DestCountry, false},
		{"OriginCountry", f.OriginCountry, false},
		{"DestCityName", f.DestCity, true},
		{"OriginCityName", f.OriginCity, true},
		{"DestAirportID", f.DestAirportID, false},
		{"OriginAirportID", f.OriginAirportID, false},
		{"DestAirport", f.DestAirport, true},
		{"OriginAirport", f.OriginAirport, true},
	}

	clauses := make([]map[string]any, 0, len(fields))
	for _, tf := range fields {
		if tf.value == nil {
			continue
		}
		key := tf.esField
        
		clauses = append(clauses, map[string]any{
			"term": map[string]any{key: *tf.value},
		})
	}
	return clauses
}

// buildCancelledFilter returns a term clause filtering on the Cancelled boolean field.
//
// Parameters:
//   - cancelled: true returns only cancelled flights, false returns only active flights
//
// Returns:
//   - []map[string]any: single-element slice with the term clause, or nil when not set
func buildCancelledFilter(f *models.FlightFilters) []map[string]any {
	if f.Cancelled == nil {
		return nil
	}
	return []map[string]any{
		{"term": map[string]any{"Cancelled": *f.Cancelled}},
	}
}

// buildDateRangeFilter returns a range clause on the timestamp field.
//
// Parameters:
//   - dateFrom: lower bound (gte), inclusive, format YYYY-MM-DD
//   - dateTo  : upper bound (lte), inclusive, format YYYY-MM-DD
//
// Returns:
//   - []map[string]any: single-element slice with the range clause, or nil when neither bound is set
func buildDateRangeFilter(f *models.FlightFilters) []map[string]any {
	if f.DateFrom == nil && f.DateTo == nil {
		return nil
	}
	r := map[string]any{}
	if f.DateFrom != nil {
		r["gte"] = f.DateFrom.Format("2006-01-02")
	}
	if f.DateTo != nil {
		r["lte"] = f.DateTo.Format("2006-01-02")
	}
	return []map[string]any{{"range": map[string]any{"timestamp": r}}}
}

// buildPriceRangeFilter returns a range clause on the AvgTicketPrice field.
//
// Parameters:
//   - priceMin: lower bound (gte), inclusive
//   - priceMax: upper bound (lte), inclusive
//
// Returns:
//   - []map[string]any: single-element slice with the range clause, or nil when neither bound is set
func buildPriceRangeFilter(f *models.FlightFilters) []map[string]any {
	if f.PriceMin == nil && f.PriceMax == nil {
		return nil
	}
	r := map[string]any{}
	if f.PriceMin != nil {
		r["gte"] = *f.PriceMin
	}
	if f.PriceMax != nil {
		r["lte"] = *f.PriceMax
	}
	return []map[string]any{{"range": map[string]any{"AvgTicketPrice": r}}}
}

// buildGeoFilter returns a geo_distance clause that restricts results to
// documents whose geo-point value for field lies within 100 km of the
// provided coordinates. The same function serves both DestLocation and
// OriginLocation by accepting the field name as a parameter.
//
// Parameters:
//   - field: Elasticsearch geo-point field name (e.g. "DestLocation", "OriginLocation")
//   - lat  : latitude in decimal degrees (-90 to 90)
//   - lon  : longitude in decimal degrees (-180 to 180)
//
// Returns:
//   - []map[string]any: single-element slice with the geo_distance clause, or nil when either coordinate is absent
func buildGeoFilter(field string, lat, lon *float64) []map[string]any {
	if lat == nil || lon == nil {
		return nil
	}
	return []map[string]any{{
		"geo_distance": map[string]any{
			"distance": "100km",
			field:      map[string]any{"lat": *lat, "lon": *lon},
		},
	}}
}