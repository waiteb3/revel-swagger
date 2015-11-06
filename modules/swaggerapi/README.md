# SwaggerAPI Revel module

## Instructions

### Installation

First grab Revel then the module and its dependency. See [revel.github.io](revel.github.io) on more details
about the revel project.
```
go get github.com/revel/cmd/revel
go get github.com/waiteb3/revel-swagger/...
go get github.com/go-swagger/go-swagger/spec
```

### Adding a SwaggerAPI route and controller

First define your Revel controller and its actions which will map to your operations
```go
// app/controllers/api.go
package controllers

import (
	"github.com/revel/revel"
)

type APIController struct {
	*revel.Controller
}

func (c APIController) ListTags(slug string) {
	// get a set of tags for an article identified by slug
}

func (c APIController) AddTag(slug, tag string) {
	// add a new tag to your article identified by slug
}
```

Then inside your Swagger specification, for each [Operation Object](http://swagger.io/specification/#operationObject)
add an `x-revel-controller-action:` property that declares which action an operation should be tied to.
```yml
swagger: "2.0"
#
# Metadata
#
paths:
  # Path Item
  /articles/{slug}:
    # Paramters for all operations on this path
    parameters:
      - in: path
        name: slug
        description: Article slug
        type: string
        required: true

    # Swagger operations for this path
    get:
      x-revel-controller-action: APIController.ListTags
      # responses &etc
    post:
      x-revel-controller-action: APIController.AddTag
      parameters:
        # This would correspond to an Admin page's submission form
        - in: formData
          name: tag
          description: Tag to be added to an article
          type: string
          required: true
      # responses &etc
```

Now you have two API endpoints that are automatically routed and parsed by Revel and SwaggerAPI.

### SwaggerAPI Configuration

All configurations go into the `conf/app.conf` file of your project.
```
# Add the module to your project
module.swaggerapi=github.com/waiteb3/revel-swagger/modules/swaggerapi
```

Add your Swagger Spec files to a location inside your project's `conf` folder.
```
# Include your files relative from the conf/ folder
# The below example would have two specifications
#
#   conf/
#       swagger.yml
#       api/
#          v1.yml
#
```

Then specify where they reside, relative to the `conf` folder, sperated by commas.
```
swaggerapi.specs=swagger.yml,api/v1.yml
```

*Note*: The built-in Swagger-UI controller relies on your module import to have the
exact name of `swaggerapi`.

#### Swagger UI config

The built-in Swagger-UI controllers are routed to `/@api`, if the basePath is specified as `/api`.
You can enable/disable the module from automatically adding the routes by setting the property to true/false.
```
swaggerapi.add-ui=true
```

Which is the effective\* equivalent of adding these for a basePath of `/api`
```
# Custom Swagger UI
GET /@api/swagger.json                       Static.Serve("conf/swagger","swagger.yml")
GET /@api/                                   Static.Serve("$MODPATH/swagger-ui/dist","index.html")
GET /@api/*filepath                          Static.Serve("$MODPATH/swagger-ui/dist")
```

\* Effective since the index.html file has to be modified to point to the swagger-spec location
and `$MODPATH` is the absolute location of this module.

### Custom UI

You can use Static.Serve (of the offical Revel module) to serve your UI.

First you'll want to turn off the built-in Swagger-UI controllers.
```
swaggerapi.add-ui=false
```

Then add these three endpoints using the Static.Serve module (order matters).
```
# Custom Swagger UI
GET /@custom/swagger.json                       Static.Serve("conf/swagger","swagger.json")
GET /@custom/                                   Static.Serve("public/swagger-ui/dist","index.html")
GET /@custom/*filepath                          Static.Serve("public/swagger-ui/dist")
```
