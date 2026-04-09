package utils

import (
	beego "github.com/beego/beego/v2/server/web"
)

// WriteError is a helper function to write a standardized JSON error response with the given HTTP status code and error message.
func WriteError(c *beego.Controller, status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]any{
		"error": message,
	}
	c.ServeJSON()
}

// WriteData is a helper function to write a standardized JSON success response with the given HTTP status code and data payload.
func WriteData(c *beego.Controller, status int, data any) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]any{
		"data": data,
	}
	c.ServeJSON()
}
