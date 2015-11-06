package controllers

import (
	"github.com/revel/revel"
	_ "github.com/waiteb3/revel-swagger/modules/swaggerapi"
)

type SwaggerAPI struct {
	*revel.Controller
}

// ServeUI renders the template for your swagger-ui
func (c SwaggerAPI) ServeUI(spec string) revel.Result {
	c.RenderArgs["SwaggerSpecURL"] = c.Request.URL.Path + spec
	return c.Render()
}
