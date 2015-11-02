package controllers

import (
	"os"
	fpath "path/filepath"

	"github.com/revel/revel"
	"github.com/waiteb3/revel-swagger/modules/swaggerapi"
)

type SwaggerAPI struct {
	*revel.Controller
}

// ServeUI renders the template for your swagger-ui
func (c SwaggerAPI) ServeUI(spec string) revel.Result {
	c.RenderArgs["AssetsURL"] = c.Request.URL.Host + c.Request.URL.Path + spec
	return c.Render()
}

// ServeAssets serves swagger-ui assets
func (c SwaggerAPI) ServeAssets(filepath string) revel.Result {
	// TODO may not be windows friendly...
	file, err := os.Open(fpath.Join(swaggerapi.ModulePath, filepath))
	if err != nil {
		return c.RenderError(err)
	}

	// TODO look into inline vs attachment
	return c.RenderFile(file, revel.Inline)
}

// Spec serves the swagger-spec file
func (c SwaggerAPI) Spec(spec string) revel.Result {
	file, err := os.Open(spec)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderFile(file, revel.Inline)
}
