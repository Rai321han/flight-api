package main

import (
	_ "flight-api/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

