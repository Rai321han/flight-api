package services

import "flight-api/models"

const flightIndex = "kibana_sample_data_flights"

func BuildSearchRequest(f *models.FlightFilters) map[string]any {
	filter := []map[string]any{}

	// (filters unchanged...)

	query := map[string]any{}
	if len(filter) > 0 {
		query["bool"] = map[string]any{"filter": filter}
	} else {
		query["match_all"] = map[string]any{}
	}

	sortField := f.SortBy
	if sortField == "" {
		sortField = "timestamp"
	}
	order := f.Order
	if order == "" {
		order = "asc"
	}

	// pagination logic
	from := (f.Page - 1) * f.Limit
	if from < 0 {
		from = 0
	}

	return map[string]any{
		"query": query,
		"sort":  []map[string]any{{sortField: map[string]any{"order": order}}},
		"size":  f.Limit,
		"from":  from,
	}
}

func IndexName() string {
	return flightIndex
}