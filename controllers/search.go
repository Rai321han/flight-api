package controllers

import (
	"errors"
	"flight-api/models"
	"flight-api/utils"
	"flight-api/validators"

	beego "github.com/beego/beego/v2/server/web"
)

// SearchController struct handles flight search requests.
// It depends on a FlightService to perform the search logic.
type SearchController struct {
	beego.Controller
	FlightSvc models.FlightService
}

// SearchFlights handles HTTP GET requests to search for flights.
// It validates query parameters, calls the FlightService to search flights,
// and returns the results in JSON format.
//
// Query Parameters:
//   - carrier, flightNum, destCountry, originCountry, destCity, originCity
//   - destAirportID, originAirportID, destAirport, originAirport
//   - dateFrom, dateTo (YYYY-MM-DD)
//   - priceMin, priceMax, cancelled
//   - destLat, destLon, originLat, originLon
//   - sortBy (timestamp|AvgTicketPrice), order (asc|desc)
//   - limit, page
//
// Responses:
//   - 200: JSON containing matching flights, total count, page, and limit.
//   - 400: Validation error for invalid query parameters.
//   - 5XX: Internal server error if search processing fails.
func (c *SearchController) SearchFlights() {

	filters, validationErrs := validators.ParseAndValidateFlightFilters(
		c.Ctx.Request.URL.Query(),
	)

	if len(validationErrs) > 0 {
		utils.WriteError(&c.Controller, 400, validationErrs[0])
		return
	}

	results, total, err := c.FlightSvc.SearchFlights(filters)
	if err != nil {
		var appErr *models.AppError
		if errors.As(err, &appErr) {
			status := appErr.Status
			if status >= 500 {
				utils.WriteError(&c.Controller, status, "internal server error")
				return
			}
			utils.WriteError(&c.Controller, status, appErr.Message)
			return
		}

		utils.WriteError(&c.Controller, 500, "internal server error")
		return
	}

	utils.WriteData(&c.Controller, 200, map[string]any{
		"flights": results,
		"total":   total,
		"page":    filters.Page,
		"limit":   filters.Limit,
	})
}
