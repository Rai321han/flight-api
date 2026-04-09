package models

// FlightService defines the contract between the controller and service layers
// for flight search operations. Placing the interface in the models package
// allows both controllers and services to depend on it without creating an
// import cycle.
type FlightService interface {
	// SearchFlights executes a flight search using the provided filters and
	// returns the matching documents, total hit count, and any error.
	//
	// Parameters:
	//   - filters: validated FlightFilters containing all active search criteria,
	//     sort preference, and pagination settings
	//
	// Returns:
	//   - []map[string]any: source documents for the current page
	//   - int: total number of documents matching the query before pagination
	//   - error: *AppError on failure, nil on success
	SearchFlights(filters *FlightFilters) ([]map[string]any, int, error)
}