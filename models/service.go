package models

// FlightService is the contract controllers depend on.
type FlightService interface {
	SearchFlights(filters *FlightFilters) ([]map[string]any, int, error)
}

// AggregationService is the contract for aggregation queries.
// Implement this interface in services/flight_aggregation_service.go
// The implementation must connect to Elasticsearch and execute aggregation queries.
type AggregationService interface {
	AggregateFlights(filters *FlightFilters) (*AggregationResult, error)
}