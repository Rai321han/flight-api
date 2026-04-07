package routers

import (
	"flight-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/flight-api/v1",
		beego.NSRouter("/search", &controllers.SearchController{}, "get:SearchFlights"),
	)

	beego.AddNamespace(ns)
}
