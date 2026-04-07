package routers

import (
	"flight-api/controllers"
	"flight-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

func Init(flightSvc models.FlightService, aggSvc models.AggregationService) {
	ns := beego.NewNamespace("/flight-api/v1",
		beego.NSRouter("/search", &controllers.SearchController{FlightSvc: flightSvc}, "get:SearchFlights"),
		beego.NSRouter("/aggregations", &controllers.AggregationController{AggSvc: aggSvc}, "get:AggregateFlights"),
	)
	beego.AddNamespace(ns)
}