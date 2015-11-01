package swaggify

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/revel/revel"
)

var ModulePath string
var AssetsPath string
var ViewsPath string
var APIs = make(map[string]*spec.Swagger)

// TODO delete or not
var IndexArgs = make(map[string]interface{})

func init() {
	_, ModulePath, _, _ = runtime.Caller(1)
	ModulePath = path.Dir(ModulePath)
	AssetsPath = filepath.Join(ModulePath, "swagger-ui", "dist")
	ViewsPath = filepath.Join(ModulePath, "app", "views")

	revel.OnAppStart(start)
}

func start() {
	// build IndexArgs for rendering index template

	// collect all of the swagger-endpoints then build the spec
	routes := make([]*revel.Route, 0)
	for _, route := range revel.MainRouter.Routes {
		if strings.ToLower(route.Action) == "swaggify.spec" {
			routes = append(routes, route)
		}
	}

	for _, route := range routes {
		// TODO test what happens on panic
		// TODO enable defaults
		// TODO complete swagger metadata and move to diff function
		// the Swagger type hoists unexported types for some annoying reason
		api := new(spec.Swagger)
		api.Swagger = "2.0"

		api.Info = new(spec.Info)
		api.Info.Version = "0.0.0"
		api.Info.Title = revel.AppName
		api.Info.Description = "Description"
		api.Info.TermsOfService = "http://swagger.io/terms/"

		api.Info.Contact = new(spec.ContactInfo)
		api.Info.Contact.Name = ""
		api.Info.Contact.Email = ""
		api.Info.Contact.URL = ""

		api.Info.License = new(spec.License)
		api.Info.License.Name = "LICENSE"
		api.Info.License.URL = "URL"

		// TODO change
		api.Host = revel.AppName
		// TODO check if https is HSTS exclusively or no for revel
		// ALSO this can be SSL terminated by proxy so this may need changing
		if revel.HttpSsl {
			api.Schemes = []string{"https"}
		} else {
			api.Schemes = []string{"http"}
		}

		// TODO configurable: JSON only for now
		api.Consumes = []string{"json"}
		api.Produces = []string{"json"}

		api.BasePath = route.FixedParams[0]

		api.Paths = buildPaths(api.BasePath)
		api.Definitions = buildDefinitions(api.BasePath)

		APIs[api.BasePath] = api
		fmt.Println(api)
	}
}

// buildPaths of a swagger spec based on the parsed revel routes and the basePath
func buildPaths(endpoint string) *spec.Paths {
	paths := new(spec.Paths)

	paths.Paths = make(map[string]spec.PathItem)

	/*
	   "paths": {
	     "/pets": {
	       "get": {
	         "description": "Returns all pets from the system that the user has access to",
	         "produces": [
	           "application/json"
	         ],
	         "responses": {
	           "200": {
	             "description": "A list of pets.",
	             "schema": {
	               "type": "array",
	               "items": {
	                 "$ref": "#/definitions/Pet"
	               }
	             }
	           }
	         }
	       }
	     }
	   },
	*/

	return paths
}

func buildDefinitions(endpoint string) spec.Definitions {
	def := spec.Definitions{}
	/*
	   "definitions": {
	     "Pet": {
	       "type": "object",
	       "required": [
	         "id",
	         "name"
	       ],
	       "properties": {
	         "id": {
	           "type": "integer",
	           "format": "int64"
	         },
	         "name": {
	           "type": "string"
	         },
	         "tag": {
	           "type": "string"
	         }
	       }
	     }
	   }
	 }
	*/
	return def
}
