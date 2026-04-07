package controllers

import (
	"flight-api/models"
	"flight-api/validators"

	beego "github.com/beego/beego/v2/server/web"
)

// ── Mock ──────────────────────────────────────────────────────────────────────

// MockAggregationService is a temporary stub so the controller compiles and
// returns realistic-shaped data while the real service is being implemented.
//
// NOTE FOR COLLEAGUE: Once FlightAggregationService is implemented in
// services/flight_aggregation_service.go, remove this mock entirely and
// wire the real service via main.go → routers.Init(aggSvc).
type MockAggregationService struct{}

func NewMockAggregationService() *MockAggregationService {
	return &MockAggregationService{}
}

// AggregateFlights returns hardcoded fake aggregation data.
//
// NOTE for you Raihan: Your real implementation must:
//  1. Call GetESClient() from services/elasticsearch_client.go
//  2. Call BuildAggregationRequest(filters) from services/aggregation_query_builder.go
//  3. Execute the query against the "kibana_sample_data_flights" index
//  4. Parse the ES response and map it into *models.AggregationResult
//  5. Return the result and any error
func (m *MockAggregationService) AggregateFlights(
	filters *models.FlightFilters,
) (*models.AggregationResult, error) {
	return &models.AggregationResult{
		FlightsPerCarrier: []models.BucketCount{
			{Key: "Kibana Airlines", Count: 3219},
			{Key: "Logstash Airways", Count: 3098},
			{Key: "ES-Air", Count: 3090},
			{Key: "JetBeats", Count: 3052},
		},
		FlightsPerDestCountry: []models.BucketCount{
			{Key: "IT", Count: 2371},
			{Key: "US", Count: 1987},
			{Key: "CN", Count: 1556},
			{Key: "CA", Count: 1234},
		},
		FlightsPerOriginCountry: []models.BucketCount{
			{Key: "IT", Count: 2344},
			{Key: "US", Count: 1902},
			{Key: "CN", Count: 1478},
			{Key: "DE", Count: 1100},
		},
		AvgPricePerCarrier: []models.BucketAvg{
			{Key: "Kibana Airlines", Count: 3219, Avg: 628.25},
			{Key: "Logstash Airways", Count: 3098, Avg: 590.10},
			{Key: "ES-Air", Count: 3090, Avg: 611.45},
			{Key: "JetBeats", Count: 3052, Avg: 598.72},
		},
		CancelledPerCarrier: []models.BucketCount{
			{Key: "Kibana Airlines", Count: 181},
			{Key: "Logstash Airways", Count: 164},
			{Key: "ES-Air", Count: 172},
			{Key: "JetBeats", Count: 159},
		},
		DelayTypeBreakdown: []models.BucketCount{
			{Key: "No Delay", Count: 10388},
			{Key: "Late Aircraft Delay", Count: 1028},
			{Key: "Security Delay", Count: 297},
			{Key: "NAS Delay", Count: 248},
			{Key: "Carrier Delay", Count: 208},
			{Key: "Weather Delay", Count: 152},
		},
		TopRoutes: []models.BucketCount{
			{Key: "ABQ → IT-0055", Count: 128},
			{Key: "XII → RDU", Count: 115},
			{Key: "DAL → VIE", Count: 109},
		},
		PriceHistogram: []models.BucketCount{
			{Key: "0", Count: 143},
			{Key: "100", Count: 892},
			{Key: "200", Count: 1784},
			{Key: "300", Count: 2201},
			{Key: "400", Count: 1998},
			{Key: "500", Count: 1643},
			{Key: "600", Count: 1204},
			{Key: "700", Count: 887},
			{Key: "800", Count: 512},
			{Key: "900", Count: 195},
		},
		FlightsPerDay: []models.BucketCount{
			{Key: "2024-01-01", Count: 143},
			{Key: "2024-01-02", Count: 155},
			{Key: "2024-01-03", Count: 148},
		},
		OverallAvgPrice: 608.25,
	}, nil
}

// ── Controller ────────────────────────────────────────────────────────────────

// AggregationController handles all aggregation-related HTTP requests.
//
// NOTE FOR you Raihan: This controller is fully implemented. Your only task is
// to implement the AggregationService interface (models/service.go) in
// services/flight_aggregation_service.go, then replace MockAggregationService
// with your real implementation in main.go.
type AggregationController struct {
	beego.Controller

	// AggSvc is injected via main.go.
	// Currently wired to MockAggregationService.
	// NOTE FOR you Raihan: Replace with services.NewFlightAggregationService()
	// in main.go once your service implementation is ready.
	AggSvc models.AggregationService
}

// AggregateFlights handles GET /flight-api/v1/aggregations
//
// It reuses ParseAndValidateFlightFilters from validators/parameters.go so the
// caller can narrow aggregations with the exact same query parameters as /search.
//
// Supported query parameters (all optional):
//   - carrier, flightNum, destCountry, originCountry, destCity, originCity
//   - destAirportID, originAirportID, destAirport, originAirport
//   - dateFrom, dateTo         (format: YYYY-MM-DD)
//   - priceMin, priceMax       (float, >= 0)
//   - cancelled                (true | false)
//   - destLat, destLon         (geo pair, both required together)
//   - originLat, originLon     (geo pair, both required together)
//   - sortBy                   (timestamp | AvgTicketPrice)
//   - order                    (asc | desc)
//   - limit                    (1–1000, default 10)
//
// Responses:
//
//	200 { "data": AggregationResult }
//	400 { "error": "invalid query parameters", "details": [...] }
//	500 { "error": "internal server error" }
func (c *AggregationController) AggregateFlights() {
	filters, validationErrs := validators.ParseAndValidateFlightFilters(c.Ctx.Request.URL.Query().Get)
	if len(validationErrs) > 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]any{
			"error":   "invalid query parameters",
			"details": validationErrs,
		}
		c.ServeJSON()
		return
	}

	result, err := c.AggSvc.AggregateFlights(filters)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]any{
			"error": "internal server error",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]any{
		"data": result,
	}
	c.ServeJSON()
}