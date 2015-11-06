package swaggerapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/revel/revel"
	"github.com/waiteb3/revel-swagger/modules/common"
)

// the extension for defined revel controller actions on swagger paths
const X_REVEL_CONTROLLER_ACTION = "x-revel-controller-action"

// init creates routes based on swagger specs defined in the conf/swagger folder
// it will create asset endpoints at /@basePath for serving the swagger-spec and UI
// which can be overridden with swaggerapi.add-ui = false
// TODO look into adding a Giles
func init() {
	revel.OnAppStart(func() {
		// TODO find out why AppRoot is empty
		//doc, err := spec.Load(filepath.Join(revel.AppRoot, "conf", "swagger", "swagger.yml"))
		//fmt.Println(revel.AppRoot, filepath.Join(revel.AppRoot, "conf", "swagger", "swagger.yml"))

		// grab the old routes to append
		oldRoutes := revel.MainRouter.Routes
		// new and emtpy router
		// NOTE WARNING calling MainRouter.Refresh will destroy the swagger routes since they do not have a routes file
		revel.MainRouter = revel.NewRouter(filepath.Join(revel.BasePath, "conf", "routes"))

		configs := revel.Config.StringDefault("swaggerapi.specs", "")
		for _, config := range strings.Split(configs, ",") {
			found := false
			for _, path := range revel.ConfPaths {
				doc, err := spec.Load(filepath.Join(path, config))
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					_fail(err)
				}

				// TODO decide if multiple of same name is allowed or not
				// if found { }

				if revel.Config.BoolDefault("swaggerapi.add-ui", true) {
					AddSwaggerUi(doc.BasePath(), config)
					// conncurently download or unzip the swagger ui assets
					// only runs once
					go common.UnzipSwaggerAssets()
				}

				AddSwaggerRoutes(doc)
				found = true
			}
			if !found {
				_fail(fmt.Errorf("Swagger configuration '%s' not found found in any of these folders:\t%v",
					config, configs))
			}
		}

		// Add oldRoutes back into the router
		for _, route := range oldRoutes {
			_fail(revel.MainRouter.Tree.Add(route.TreePath, route))
			revel.MainRouter.Routes = append(revel.MainRouter.Routes, route)
		}
	})
}

func AddSwaggerRoutes(doc *spec.Document) {
	for path, pathItem := range doc.AllPaths() {
		path = doc.BasePath() + deCurly(path)

		// TODO may be something usable here?
		// if extension, ok := pathItem.Extensions.GetString("x-revel-controller"); ok { }

		if pathItem.Head != nil {
			if action, ok := pathItem.Head.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "HEAD", path))
			}
		}

		if pathItem.Get != nil {
			if action, ok := pathItem.Get.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "GET", path))
			}
		}

		if pathItem.Post != nil {
			if action, ok := pathItem.Post.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "POST", path))
			}
		}

		if pathItem.Put != nil {
			if action, ok := pathItem.Put.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "PUT", path))
			}
		}

		if pathItem.Delete != nil {
			if action, ok := pathItem.Delete.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "DELETE", path))
			}
		}

		if pathItem.Patch != nil {
			if action, ok := pathItem.Patch.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_fail(addRoute(action, "PATCH", path))
			}
		}

		if pathItem.Options != nil {
			if action, ok := pathItem.Options.Extensions.GetString("x-revel-controller-action"); ok {
				_fail(addRoute(action, "OPTIONS", path))
			}
		}
	}
}

// AddSwaggerUI inserts the routes for serving UIs at /@basePath/
func AddSwaggerUi(basePath, filename string) {
	basePath = insertAtSymbol(basePath)

	// GET /@{basePath}/{spec}   Static.Serve("{projPath}","conf/{spec}")
	_fail(addRoute("Static.Serve", "GET", basePath+"/"+filename,
		revel.BasePath, filepath.Join("conf", filename)))

	// GET /@{basePath}/	     SwaggerAPI.ServeUI("{spec}")
	_fail(addRoute("SwaggerAPI.ServeUI", "GET", basePath+"/",
		filename))

	// GET /@{basePath}/{spec}   Static.Serve("{modPath}/swagger-ui/dist")
	_fail(addRoute("Static.Serve", "GET", basePath+"/*filepath",
		common.SwaggerAssetsDir))
}

// adds a new route to the MainRouter
func addRoute(action, method, path string, params ...string) error {
	i := strings.LastIndex(action, ".")
	if i < 0 {
		return fmt.Errorf("Invalid parsing of x-revel-controller-action '%s' for HTTP %s at '%s",
			action, method, path)
	}

	ControllerName, MethodName := action[:i], action[i+1:]
	treePath := "/" + method + path
	route := &revel.Route{
		Action:         action,
		ControllerName: ControllerName,
		Path:           path,
		TreePath:       treePath,
		Method:         method,
		MethodName:     MethodName,
		FixedParams:    params,
	}

	err := revel.MainRouter.Tree.Add(treePath, route)
	if err != nil {
		return err
	}

	revel.MainRouter.Routes = append(revel.MainRouter.Routes, route)
	return nil
}

// deCurly replaces swagger {curlyparams} with revel's :colonparams
func deCurly(path string) string {
	path = strings.Map(func(r rune) rune {
		if r == '{' {
			return ':'
		}
		return r
	}, path)
	path = strings.Replace(path, "}", "", -1)
	return path
}

// insertAtSymbol returns the basePath with an @ inseted before the first /
// if basePath is an empty string, it returns /@
func insertAtSymbol(basePath string) string {
	if basePath == "" {
		return "/@"
	}
	return basePath[:1] + "@" + basePath[1:]
}

// TODO get rid of what is basically a goto fail?
func _fail(err error) {
	if err != nil {
		revel.ERROR.Fatalln(err)
	}
}
