package models

import "time"

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