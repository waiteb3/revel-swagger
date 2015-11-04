# SwaggerAPI Revel module

## Instructions

First grab the module
```
go get github.com/waiteb3/revel-swagger/modules/swaggerapi/...
```

### SwaggerAPI Configuration

All configurations go into the `conf/app.conf` file of your project.
```
# Add the module to your project
module.swaggerapi=github.com/waiteb3/revel-swagger/modules/swaggerapi
```

Add your Swagger Spec files to a location inside your project's `conf` folder
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

Enable/disable the built-in Swagger-UI serving controllers/
```
swaggerapi.add-ui=true
# false disables them
```

The built-int Swagger-UI controllers are routed to `/@api`, if the basePath is specified as `/api`

### Custom UI

You can use Static.Serve (the offical Revel module) to serve your UI.

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
