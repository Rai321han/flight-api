package main

import (
	"flight-api/routers"
	"flight-api/services"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	flightSvc := services.NewFlightSearchService()
	routers.Init(flightSvc)
	beego.Run()
}