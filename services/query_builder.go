package services

import "flight-api/models"

const flightIndex = "kibana_sample_data_flights"

var flightSourceFields = []string{
	"FlightNum", "Carrier", "FlightDelay", "FlightDelayMin", "FlightDelayType",
	"OriginCityName", "OriginCountry", "OriginAirportID", "OriginAirport",
	"OriginLocation",
	"DestCityName", "DestCountry", "DestAirportID", "DestAirport",
	"DestLocation",
	"AvgTicketPrice", "Cancelled", "timestamp",
}

type termField struct {
	esField string
	value   *string
	keyword bool
}

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
		"query":    query,
		"sort":     []map[string]any{{sortField: map[string]any{"order": order}}},
		"size":     f.Limit,
		"from":     from,
		"_source":  flightSourceFields,
	}
}

func IndexName() string {
	return flightIndex
}

// buildFilter assembles the bool.filter slice from all active FlightFilters.
// Complexity is kept low by delegating each concern to a dedicated builder.
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

func buildTermFilters(f *models.FlightFilters) []map[string]any {
	fields := []termField{
		{"Carrier", f.Carrier, true},
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
		if tf.keyword {
			key += ".keyword"
		}
		clauses = append(clauses, map[string]any{
			"term": map[string]any{key: *tf.value},
		})
	}
	return clauses
}

func buildCancelledFilter(f *models.FlightFilters) []map[string]any {
	if f.Cancelled == nil {
		return nil
	}
	return []map[string]any{
		{"term": map[string]any{"Cancelled": *f.Cancelled}},
	}
}

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