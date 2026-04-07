package controllers

import (
	"flight-api/models"
	"flight-api/utils"
	"flight-api/validators"

	beego "github.com/beego/beego/v2/server/web"
)

// allowedSortFields restricts sortBy to known ES fields, preventing injection.
var allowedSortFields = map[string]bool{
	"AvgTicketPrice": true,
	"timestamp":      true,
}

type SearchController struct {
	beego.Controller
	FlightSvc models.FlightService
}

// SearchFlights handles GET /flight-api/v1/search
func (c *SearchController) SearchFlights() {
	filters, validationErrs := validators.ParseAndValidateFlightFilters(c.Ctx.Request.URL.Query().Get)
	if len(validationErrs) > 0 {
		utils.WriteError(&c.Controller, 400, "invalid query parameters", validationErrs)
		return
	}

	results, _, err := c.FlightSvc.SearchFlights(filters)
	if err != nil {
		utils.WriteError(&c.Controller, 500, "internal server error", nil)
		return
	}
	utils.WriteData(&c.Controller, 200, results)
}
