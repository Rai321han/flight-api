package utils

import (
	beego "github.com/beego/beego/v2/server/web"
)

type APIError struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func WriteError(c *beego.Controller, status int, message string, details []string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]any{
		"error": APIError{
			Message: message,
			Details: details,
		},
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
