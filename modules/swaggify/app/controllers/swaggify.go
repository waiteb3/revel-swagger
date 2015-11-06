package contollers

import (
	"github.com/revel/modules/static/app/controllers"
	"github.com/revel/revel"
	"github.com/waiteb3/revel-swagger/modules/common"
	"github.com/waiteb3/revel-swagger/modules/swaggify"
)

type Swaggify struct {
	*revel.Controller
}

func (c Swaggify) ServeUI(basePath string) revel.Result {
	c.RenderArgs["SwaggerSpecURL"] = c.Request.URL.Path + "swagger.json"
	return c.Render()
}

var assetsFixedParams map[string][]string

func init() {
	revel.OnAppStart(func() {
		assetsFixedParams = map[string][]string{
			"prefix": []string{common.SwaggerAssetsDir},
		}
	})
}

func (c Swaggify) ServeAssets(filepath string) revel.Result {
	cast := controllers.Static{c.Controller}
	// SEE https://github.com/revel/modules/blob/master/static/app/controllers/static.go#L60
	c.Params.Fixed = assetsFixedParams
	return cast.Serve(common.SwaggerAssetsDir, filepath)
}

func (c Swaggify) Spec(basePath string) revel.Result {
	if api := swaggify.APIs[basePath]; api != nil {
		// if this is accessed from two different hostnames it
		// risks a torn read/write
		// api.Host = c.Request.URL.Host
		return c.RenderJson(api)
	} else {
		return c.NotFound("No matching API at endpoint %s", basePath)
	}
}
