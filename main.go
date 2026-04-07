package main

import (
	"flight-api/controllers"
	"flight-api/routers"
	"flight-api/services"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// Live Elasticsearch service — fully implemented.
	flightSvc := services.NewFlightSearchService()

	// NOTE FOR you Raihan: MockAggregationService is a temporary stub.
	// Once you have implemented FlightAggregationService in
	// services/flight_aggregation_service.go, replace the line below with:
	//
	//   aggSvc := services.NewFlightAggregationService()
	//
	// Then delete controllers.NewMockAggregationService() and the
	// MockAggregationService struct from controllers/aggregation.go entirely.
	aggSvc := controllers.NewMockAggregationService()

	routers.Init(flightSvc, aggSvc)
	beego.Run()
}