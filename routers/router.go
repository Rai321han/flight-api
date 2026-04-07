package routers

import (
	"flight-api/controllers"
	"flight-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

func Init(flightSvc models.FlightService) {
	ns := beego.NewNamespace("/flight-api/v1",
		beego.NSRouter("/search", &controllers.SearchController{FlightSvc: flightSvc}, "get:SearchFlights"),
	)
	beego.AddNamespace(ns)
}