# Revel-Swagger
Two drop in Revel modules for Swagger integration

## Modules quickstart

First grab the modules
```
go get github.com/waiteb3/revel-swagger/modules/...
```

### SwaggerAPI

#### Introduction

SwaggerAPI is a way to add routing based on a Swagger Specification file into your Revel project.

#### Quickstart

Add the module and spec locations to `app.conf`
```
# Add the module
module.swaggerapi=github.com/waiteb3/revel-swagger/modules/swaggerapi
# Include your spec relative from the conf/ folder
swaggerapi.specs=swagger.yml
```

Add your Swagger-Spec file to the `swagger` folder
```
ls conf
# app.conf
# routes
# swagger.yml
```

Done! This will generate the routes on start-up and will begin to serve the Swagger-UI assets at `/@{basePath}`

See the [example project](examples/swaggerapi) and a complete introduction at module's [README.md](modules/swaggerapi/README.md).

If you wish to serve your own Swagger-UI distribution, see the section on using [Static.Serve](modules/swaggerapi/README.md#custom-ui).

**Note**: Currently you have to parse the contents of a c.Request.Body yourself. See [Item.Create](examples/swaggerapi/app/controllers/items.go#L51-L62) for an example.

### (Functional WIP) Swaggify

#### Introduction

Swaggify is a module that generates a Swagger Specification best-guess based on the source-code of your controllers (and hopefully models as well).

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


### (Planned/Maybe) Editor

Add the module to `conf/app.conf`
```
module.swaggereditor=github.com/waiteb3/revel-swagger/modules/editor
```

Drop the Swagger spec editor endpoint into `conf/routes`
```
GET     /@api/editor                          SwaggerEditor.Serve("conf", "swagger.yml")
```

#### Example

## Contributing
Contributions are welcome in the form of issues or pull requests. Minor changes are fine,
but for major changes, please open an empty pull request for discussion or raise an issue first.
The last thing I want is for you to make a few hundred changes and then get told no.
Suggestions and civil complaints are also fine in the form of an issue.
