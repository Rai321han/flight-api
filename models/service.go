package models

// FlightService is the contract controllers depend on.
// Defined in models so both controllers and services can reference it
// without creating an import cycle.
type FlightService interface {
	SearchFlights(filters *FlightFilters) ([]map[string]any, int, error)
}