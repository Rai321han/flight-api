package controllers

import (
	"flight-api/models"
	"flight-api/utils"
	"flight-api/validators"
	"github.com/beego/beego/v2/server/web"
)

// SearchController handles flight search requests
type SearchController struct {
	web.Controller
	FlightSvc models.FlightService
}

// SearchFlights handles GET /flight-api/v1/search
func (c *SearchController) SearchFlights() {
	// Parse query parameters
	filters, errs := validators.ParseAndValidateFlightFilters(c.Ctx.Request.URL.Query().Get)
	if len(errs) > 0 {
		utils.WriteError(&c.Controller, 400, "Invalid query parameters")
		return
	}

	// Call the FlightService to search flights
	results, total, err := c.FlightSvc.SearchFlights(filters)
	if err != nil {
		utils.WriteError(&c.Controller, 500, "Internal server error")
		return
	}

	// Return paginated response
	utils.WriteData(&c.Controller, 200, map[string]any{
		"flights": results,
		"total":   total,
		"limit":   filters.Limit,
	})
}