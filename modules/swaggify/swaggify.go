package swaggify

import (
	"fmt"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/revel/cmd/harness"
	"github.com/revel/revel"
	"github.com/waiteb3/revel-swagger/modules/common"
)

var APIs = make(map[string]*spec.Swagger)

// TODO try to add manual defintion override identified by something in the comments
// TODO interceptor for getting statistics on endpoints
// TODO detect what types were produced
var ContentTypes = []string{"application/json"}

func init() {
	revel.OnAppStart(func() {
		go common.UnzipSwaggerAssets()
		// build IndexArgs for rendering index template

		// collect all of the swagger-endpoints then build the spec
		routes := make([]*revel.Route, 0)
		for _, route := range revel.MainRouter.Routes {
			if strings.ToLower(route.Action) == "swaggify.spec" {
				routes = append(routes, route)
			}
		}

		for _, route := range routes {
			// TODO test what cases cause bounds panic
			// Don't duplicate building API specs
			if _, exists := APIs[route.FixedParams[0]]; exists {
				continue
			}

			APIs[route.FixedParams[0]] = newSpec(route.FixedParams[0])
			fmt.Println(APIs[route.FixedParams[0]])
		}
	})
}

func newSpec(endpoint string) *spec.Swagger {
	// TODO enable defaults
	// TODO complete swagger metadata and move to diff function
	api := new(spec.Swagger)
	api.Swagger = "2.0"
	api.BasePath = endpoint

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

	// Now is added on requests to Swaggify.Spec

	// TODO check if https is HSTS exclusively or no for revel
	// TODO ALSO this can be SSL terminated by proxy so this may need changing
	if revel.HttpSsl {
		api.Schemes = []string{"https"}
	} else {
		api.Schemes = []string{"http"}
	}

	api.Consumes = ContentTypes
	api.Produces = ContentTypes

	// the Swagger type hoists unexported types for some annoying reason
	api.Paths = buildPaths(api.BasePath)
	api.Definitions = buildDefinitions(api.BasePath)

	return api
}

// buildPaths of a swagger spec based on the parsed revel routes and the basePath
func buildPaths(endpoint string) *spec.Paths {
	paths := new(spec.Paths)

	paths.Paths = make(map[string]spec.PathItem)

	for _, route := range revel.MainRouter.Routes {
		// swagger uses curly braces and prepends the basePath
		if strings.HasPrefix(route.Path, endpoint) {
			swaggerPath := curlify(endpoint, route.Path)
			// match only routes in the endpoint
			fmt.Println()
			fmt.Println("endpoint", route)

			// missing values returns an empty struct
			path := paths.Paths[swaggerPath]
			// TODO MAYBE path.Extensions
			switch strings.ToUpper(route.Method) {
			case "HEAD":
				path.Head = buildOperation(route)
			case "GET":
				path.Get = buildOperation(route)
			case "POST":
				path.Post = buildOperation(route)
			case "PUT":
				path.Put = buildOperation(route)
			case "DELETE":
				path.Delete = buildOperation(route)
			case "PATCH":
				path.Patch = buildOperation(route)
			case "OPTIONS":
				path.Options = buildOperation(route)
			case "*":
				path.Head = buildOperation(route)
				path.Get = path.Head
				path.Post = path.Head
				path.Put = path.Head
				path.Delete = path.Head
				path.Patch = path.Head
				path.Options = path.Head
			default:
				panic("Invalid Request type ?? " + route.Method)
			}
			paths.Paths[swaggerPath] = path
		}
	}

	return paths
}

// build an operation object based on the route information
func buildOperation(route *revel.Route) *spec.Operation {
	var (
		typeInfo   *harness.TypeInfo
		methodSpec *harness.MethodSpec
	)

	info, rerr := harness.ProcessSource(revel.CodePaths)
	if rerr != nil {
		panic(rerr) // TODO EMPTY PANIC
	}

	// get the TypeInfo and MethodSpec for this route
	for _, cinfo := range info.ControllerSpecs() {
		typeInfo = cinfo // TODO move inside if (get around compiler complaint)
		if route.ControllerName == typeInfo.StructName {
			for _, spec := range cinfo.MethodSpecs {
				if route.MethodName == spec.Name {
					methodSpec = spec
					break
				}
			}
			break
		}
	}

	op := new(spec.Operation)
	// TODO op.Description
	// this will probably require extending harness.ProcessSource to parse comments
	op.Consumes = ContentTypes
	op.Produces = ContentTypes
	op.AddExtension("x-revel-controller-action", route.Action)

	op.Tags = []string{trimControllerName(route.ControllerName)}

	for i, arg := range methodSpec.Args {
		// skip over fixed paramters that match up to the arguments
		if i < len(route.FixedParams) {
			continue
		}
		var param spec.Parameter
		param.Name = arg.Name
		param.Type = arg.TypeExpr.Expr

		// TODO review
		// TODO: better path vs query vs body vs multipart
		count := strings.Count(route.Path, ":") + strings.Count(route.Path, "*")
		if i < count {
			param.In = "path"
		} else {
			param.In = "formData"
		}
		op.Parameters = append(op.Parameters, param)
	}

	// TODO RenderCalls
	// fmt.Printf("route:       %#v\n", route)
	// fmt.Printf("typeInfo:    %#v\n", typeInfo)
	// fmt.Printf("methodSpec: %#v\n", methodSpec)
	// for _, call := range methodSpec.RenderCalls {
	// 	fmt.Printf("\tcall: %#v\n", call)
	// }

	/*
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
	*/

	return op
}

// Trim off "Controller(s)" from the ControllerName to use as a tag
var trimControllerName = func(name string) string {
	if strings.HasSuffix(strings.ToLower(name), "controller") {
		return name[:len(name)-len("controller")]
	}
	if strings.HasSuffix(strings.ToLower(name), "controllers") {
		return name[:len(name)-len("controllers")]
	}
	return name
}

// TODO Parse the models and contollers for potential definitions?
func buildDefinitions(endpoint string) spec.Definitions {
	def := spec.Definitions{}

	// TODO
	// info, err := harness.ProcessSource(revel.CodePaths...)
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

func curlify(endpoint, path string) string {
	path = strings.TrimPrefix(path, endpoint)
	parts := strings.SplitAfter(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
			parts[i] = "{" + part[1:] + "}"
		}
	}
	return strings.Join(parts, "")
}
