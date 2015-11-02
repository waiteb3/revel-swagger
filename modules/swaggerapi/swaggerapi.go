package swaggerapi

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/revel/revel"
)

var ModulePath string

// the extension for defined revel controller actions on swagger paths
const X_REVEL_CONTROLLER_ACTION = "x-revel-controller-action"

// init creates routes based on swagger specs defined in the conf/swagger folder
// it will create asset endpoints at /@basePath for serving the swagger-spec and UI
// which can be overridden with swaggerapi.addui = false
func init() {
	_, ModulePath, _, _ = runtime.Caller(1)
	ModulePath = filepath.Join(path.Dir(ModulePath), "swagger-ui", "dist")

	revel.OnAppStart(func() {
		// TODO look into alternatives rather than than using the config (maybe)

		// TODO find out why AppRoot is empty
		//doc, err := spec.Load(filepath.Join(revel.AppRoot, "conf", "swagger", "swagger.yml"))
		//fmt.Println(revel.AppRoot, filepath.Join(revel.AppRoot, "conf", "swagger", "swagger.yml"))

		configs := revel.Config.StringDefault("swaggerapi.specs", "")
		for _, config := range strings.Split(configs, ",") {
			found := false
			for _, path := range revel.ConfPaths {
				doc, err := spec.Load(filepath.Join(path, config))
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					_panic(err)
				}

				// TODO decide if multiple of same name is allowed or not
				// if found { }

				if revel.Config.BoolDefault("swaggerapi.no-ui", true) {
					AddSwaggerUi(doc.BasePath(), config)
				}

				AddSwaggerRoutes(doc)
				found = true
			}
			if !found {
				_panic(fmt.Errorf("Swagger configuration '%s' not found found in any folder of these folders:\t%v",
					config, configs))
			}
		}
	})
}

func AddSwaggerRoutes(doc *spec.Document) {
	for path, pathItem := range doc.AllPaths() {
		path = doc.BasePath() + deCurly(path)

		// TODO may be something usable here?
		// if extension, ok := pathItem.Extensions.GetString("x-revel-controller"); ok { }

		// TODO think on how to clena this up some more
		if pathItem.Head != nil {
			if action, ok := pathItem.Head.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "HEAD", path))
			}
		}

		if pathItem.Get != nil {
			if action, ok := pathItem.Get.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "GET", path))
			}
		}

		if pathItem.Post != nil {
			if action, ok := pathItem.Post.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "POST", path))
			}
		}

		if pathItem.Put != nil {
			if action, ok := pathItem.Put.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "PUT", path))
			}
		}

		if pathItem.Delete != nil {
			if action, ok := pathItem.Delete.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "DELETE", path))
			}
		}

		if pathItem.Patch != nil {
			if action, ok := pathItem.Patch.Extensions.GetString(X_REVEL_CONTROLLER_ACTION); ok {
				_panic(addRoute(action, "PATCH", path))
			}
		}

		if pathItem.Options != nil {
			if action, ok := pathItem.Options.Extensions.GetString("x-revel-controller-action"); ok {
				_panic(addRoute(action, "OPTIONS", path))
			}
		}
	}
}

// AddSwaggerUI inserts the routes for serving UIs at /@basePath/
func AddSwaggerUi(basePath, filename string) {
	basePath = insertAtSymbol(basePath)

	// TODO don't stop down already defined /@path's
	// leaf, _ := revel.MainRouter.Tree.Find(basePath + "/" + filename)
	_panic(addRoute("SwaggerAPI.Spec", "GET", basePath+"/"+filename,
		filepath.Join("conf", filename)))

	_panic(addRoute("SwaggerAPI.ServeUI", "GET", basePath+"/",
		filename))

	_panic(addRoute("SwaggerAPI.ServeAssets", "GET", basePath+"/*filepath"))
}

// adds a new route to the MainRouter
func addRoute(action, method, path string, params ...string) error {
	i := strings.LastIndex(action, ".")
	if i < 0 {
		return fmt.Errorf("Invalid parsing of x-revel-controller-action '%s' for HTTP %s at '%s",
			action, method, path)
	}

	ControllerName, MethodName := action[:i], action[i+1:]
	return revel.MainRouter.Tree.Add("/"+method+path, &revel.Route{
		Action:         action,
		ControllerName: ControllerName,
		Path:           path,
		TreePath:       "/" + method + path,
		Method:         method,
		MethodName:     MethodName,
		FixedParams:    params,
	})
}

// deCurly replaces swagger {curlyparams} with revel's :colonparams
func deCurly(path string) string {
	// TODO gotta get this done
	path = strings.Map(func(r rune) rune {
		if r == '{' {
			return ':'
		}
		return r
	}, path)
	path = strings.Replace(path, "}", "", -1)
	fmt.Println("deCurly", path)
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

// basically goto fail =)
func _panic(err error) {
	if err != nil {
		panic(err)
	}
}
