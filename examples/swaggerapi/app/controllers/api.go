package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Endpoint() revel.Result {
	return c.RenderJson(struct {
		Message string `json:"message"`
	}{"Swagger routing!"})
}
