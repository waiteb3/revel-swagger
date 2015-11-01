# Revel-Swagger
Two drop in Revel modules for Swagger integration

## Modules quickstart

### Swaggify

Add the module to `conf/app.conf`
```
module.swaggify=github.com/waiteb3/revel-swagger/modules/swaggify
```

Drop the endpoints into `conf/routes`
```
# This will serve the generated swagger spec built from from any route with
# a basePath of "/api" and spec accessible at "/@api/swagger.json"
GET     /@api/swagger.json                      Swaggify.Spec("/api")
# Serve the default swagger-ui based on the prefix ("/api")
GET     /@api/	                                Swaggify.ServeUI("/api")
GET     /@api/*filepath                         Swaggify.ServeAssets("/api")
```

#### Example

### (Planned) Router

Add the module to `conf/app.conf`
```
module.swaggerapi=github.com/waiteb3/revel-swagger/modules/swaggerapi
```

Drop the endpoints into `conf/routes`
```
# This will create a router based on the swagger.(yml|json)
*       /api                                  SwaggerAPI.Endpoint("conf", "swagger.yml")
```

#### Example

### (Planned) Editor

Add the module to `conf/app.conf`
```
module.swaggereditor=github.com/waiteb3/revel-swagger/modules/editor
```

Drop the Swagger spec editor endpoint into `conf/routes`
```
GET     /@api/editor                          SwaggerEditor.Serve("conf", "swagger.yml")
```

#### Example

