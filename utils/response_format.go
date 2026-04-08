package utils

import (
	beego "github.com/beego/beego/v2/server/web"
)

func WriteError(c *beego.Controller, status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]any{
		"error": message,
	}
	c.ServeJSON()
}

func WriteData(c *beego.Controller, status int, data any) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]any{
		"data": data,
	}
	c.ServeJSON()
}
