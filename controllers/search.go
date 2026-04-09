package controllers

import (
	"errors"
	"flight-api/models"
	"flight-api/utils"
	"flight-api/validators"

	beego "github.com/beego/beego/v2/server/web"
)

type SearchController struct {
	beego.Controller
	FlightSvc models.FlightService
}

// SearchFlights handles GET /search endpoint to search for flights based on various filters.
// It validates query parameters, calls the FlightService to perform the search, and returns the results in a paginated format.
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
