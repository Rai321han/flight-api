package services

import "flight-api/models"

// FlightService is the contract controllers depend on.
type FlightService interface {
	SearchFlights(filters *models.FlightFilters) ([]map[string]any, int, error)
}