package controllers

import (
	"errors"
	"flight-api/models"
	"flight-api/utils"
	"flight-api/validators"
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

type SearchController struct {
	beego.Controller
	FlightSvc models.FlightService
}

// SearchFlights handles GET /flight-api/v1/search
func (c *SearchController) SearchFlights() {
	filters, err := validators.ParseAndValidateFlightFilters(c.Ctx.Request.URL.Query().Get)
	if err != nil {
		c.handleValidationError(err)
		return
	}

	results, total, err := c.FlightSvc.SearchFlights(filters)
	if err != nil {
		c.handleServiceError(err)
		return
	}

	utils.WriteData(&c.Controller, 200, map[string]any{
		"flights": results,
		"total":   total,
	})
}

func (c *SearchController) handleValidationError(err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		utils.WriteError(&c.Controller, appErr.Status, appErr.Message)
		return
	}

	log.Printf("[SearchController] validation error: %v", err)
	utils.WriteError(&c.Controller, 400, "invalid query parameters")
}

func (c *SearchController) handleServiceError(err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		if appErr.Cause != nil {
			log.Printf("[SearchController] search failed: %v", appErr.Cause)
		}

		message := appErr.Message
		if appErr.Status >= 500 {
			message = "internal server error"
			if appErr.Cause != nil {
				log.Printf("[SearchController] %d error - actual: %s: %v", appErr.Status, appErr.Message, appErr.Cause)
			} else {
				log.Printf("[SearchController] %d error - actual: %s", appErr.Status, appErr.Message)
			}
		}

		utils.WriteError(&c.Controller, appErr.Status, message)
		return
	}

	log.Printf("[SearchController] unexpected error: %v", err)
	utils.WriteError(&c.Controller, 500, "internal server error")
}
